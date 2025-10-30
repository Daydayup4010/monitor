package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	"uu/config"
	"uu/models"
	"uu/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type WechatLoginRequest struct {
	Code string `json:"code" binding:"required"`
}

type WechatSessionResponse struct {
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionID    string `json:"unionid"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

type BindEmailRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Code     string `json:"code" binding:"required,len=6"`
	Password string `json:"password" binding:"required,min=6"`
}

func WechatLogin(c *gin.Context) {
	var req WechatLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": utils.ErrCodeInvalidParams,
			"msg":  utils.ErrorMessage(utils.ErrCodeInvalidParams),
		})
		return
	}

	// 1. Get wechat openid from code
	openID, err := getWechatOpenID(req.Code)
	if err != nil {
		config.Log.Errorf("get wechat openid error: %v", err)
		c.JSON(http.StatusOK, gin.H{
			"code": utils.ErrCodeWechatLogin,
			"msg":  utils.ErrorMessage(utils.ErrCodeWechatLogin),
		})
		return
	}

	// 2. Query user by openid
	user := models.QueryUserByOpenID(openID)

	if user == nil {
		// 3a. New user: auto create account
		userId := uuid.New()
		user = &models.User{
			ID:           userId,
			UserName:     "wechat_user_" + openID[len(openID)-6:],
			WechatOpenID: &openID,
			Role:         models.RoleNormal,
			LastLogin:    time.Now(),
		}

		code := models.CreateDefaultSetting(userId.String())
		if code != utils.SUCCESS {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": code,
				"msg":  utils.ErrorMessage(code),
			})
			return
		}

		code = models.CreateWechatUser(user)
		if code != utils.SUCCESS {
			c.JSON(http.StatusOK, gin.H{
				"code": code,
				"msg":  utils.ErrorMessage(code),
			})
			return
		}
		config.Log.Infof("new wechat user registered: %s", user.ID.String())
	} else {
		// 3b. Existing user: update last login time
		user.LastLogin = time.Now()
		models.UpdateUserLastLogin(user)
		config.Log.Infof("wechat user login: %s", user.ID.String())
	}

	// 4. Generate JWT token
	token, err := utils.GenerateJWT(user.ID, user.UserName, user.Role, user.VipExpiry, user.Email)
	if err != nil {
		config.Log.Errorf("generate token error: %v", err)
		c.JSON(http.StatusOK, gin.H{
			"code": utils.ErrCodeTokenGenerate,
			"msg":  utils.ErrorMessage(utils.ErrCodeTokenGenerate),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":  utils.SUCCESS,
		"token": token,
		"data": gin.H{
			"id":         user.ID.String(),
			"username":   user.UserName,
			"email":      user.Email,
			"role":       user.Role,
			"vip_expiry": user.VipExpiry,
			"has_email":  user.Email != "",                                     // Has email bound
			"has_wechat": user.WechatOpenID != nil && *user.WechatOpenID != "", // Has wechat bound
		},
		"msg": utils.ErrorMessage(utils.SUCCESS),
	})
}

func BindEmail(c *gin.Context) {
	var req BindEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": utils.ErrCodeInvalidParams,
			"msg":  utils.ErrorMessage(utils.ErrCodeInvalidParams),
		})
		return
	}

	userID := getUserIdFromContext(c)

	// 1. Verify email code
	valid, code := models.VerifyEmailCode(req.Email, req.Code, c)
	if !valid {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  utils.ErrorMessage(code),
		})
		return
	}

	// 2. Check if email already exists
	if models.IfExistEmail(req.Email) {
		c.JSON(http.StatusOK, gin.H{
			"code": utils.ErrCodeEmailTaken,
			"msg":  utils.ErrorMessage(utils.ErrCodeEmailTaken),
		})
		return
	}

	// 3. Bind email and password
	hashedPassword := models.ScryptPw(req.Password)
	updates := map[string]interface{}{
		"email":    req.Email,
		"password": hashedPassword,
	}

	err := config.DB.Model(&models.User{}).Where("id = ?", userID).Updates(updates).Error
	if err != nil {
		config.Log.Errorf("bind email error: %v", err)
		c.JSON(http.StatusOK, gin.H{
			"code": utils.ErrCodeUpdateUser,
			"msg":  utils.ErrorMessage(utils.ErrCodeUpdateUser),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": utils.SUCCESS,
		"msg":  utils.ErrorMessage(utils.SUCCESS),
	})
}

func BindWechat(c *gin.Context) {
	var req WechatLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": utils.ErrCodeInvalidParams,
			"msg":  utils.ErrorMessage(utils.ErrCodeInvalidParams),
		})
		return
	}

	userID := getUserIdFromContext(c)

	// 1. Get wechat openid
	openID, err := getWechatOpenID(req.Code)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": utils.ErrCodeWechatLogin,
			"msg":  utils.ErrorMessage(utils.ErrCodeWechatLogin),
		})
		return
	}

	// 2. Check if this openid is already bound to another user
	existUser := models.QueryUserByOpenID(openID)
	if existUser != nil && existUser.ID.String() != userID {
		c.JSON(http.StatusOK, gin.H{
			"code": utils.ErrCodeWechatBindFailed,
			"msg":  utils.ErrorMessage(utils.ErrCodeWechatBindFailed),
		})
		return
	}

	// 3. Bind wechat
	updates := map[string]interface{}{
		"wechat_openid": openID,
	}

	err = config.DB.Model(&models.User{}).Where("id = ?", userID).Updates(updates).Error
	if err != nil {
		config.Log.Errorf("bind wechat error: %v", err)
		c.JSON(http.StatusOK, gin.H{
			"code": utils.ErrCodeUpdateUser,
			"msg":  utils.ErrorMessage(utils.ErrCodeUpdateUser),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": utils.SUCCESS,
		"msg":  utils.ErrorMessage(utils.SUCCESS),
	})
}

// Call wechat API to get openid
func getWechatOpenID(code string) (string, error) {
	// Read from config file
	appID := config.CONFIG.Wechat.AppID
	appSecret := config.CONFIG.Wechat.AppSecret

	url := fmt.Sprintf(
		"https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
		appID, appSecret, code,
	)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			config.Log.Errorf("Get wechat openid error")
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result WechatSessionResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	if result.ErrCode != 0 {
		return "", fmt.Errorf("wechat api error: %s", result.ErrMsg)
	}

	return result.OpenID, nil
}
