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

// 合并账号请求
type MergeAccountRequest struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required,len=6"`
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
		// 3a. New user: auto create account with 2 days VIP trial
		userId := uuid.New()
		vipExpiry := time.Now().Add(2 * 24 * time.Hour) // 新用户赠送2天VIP试用
		user = &models.User{
			ID:           userId,
			UserName:     "wechat_user_" + openID[len(openID)-6:],
			WechatOpenID: &openID,
			Role:         models.RoleVip, // 新用户默认VIP
			VipExpiry:    vipExpiry,
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
		config.Log.Infof("new wechat user registered: %s, granted 2 days VIP trial until %v", user.ID.String(), vipExpiry)
	} else {
		// 3b. Existing user: update last login time
		user.LastLogin = time.Now()
		models.UpdateUserLastLogin(user)
		config.Log.Infof("wechat user login: %s", user.ID.String())
	}

	// 4. Generate token version and store to Redis (按客户端类型区分，小程序端)
	tokenVersion := models.GenerateTokenVersion()
	if err := models.SetTokenVersion(c.Request.Context(), user.ID, models.ClientTypeMiniprogram, tokenVersion); err != nil {
		config.Log.Errorf("set token version error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": utils.ErrCodeTokenGenerate,
			"msg":  "登录失败，请重试",
		})
		return
	}

	// 5. Generate JWT token
	token, err := utils.GenerateJWT(user.ID, user.UserName, user.Role, user.VipExpiry, user.Email, tokenVersion, models.ClientTypeMiniprogram)
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

	// 2. Check if email already exists
	if models.IfExistEmail(req.Email) {
		// 邮箱已存在，提示需要合并账号
		c.JSON(http.StatusOK, gin.H{
			"code": utils.ErrCodeEmailExistsNeedMerge,
			"msg":  utils.ErrorMessage(utils.ErrCodeEmailExistsNeedMerge),
		})
		return
	}

	// 1. Verify email code
	valid, code := models.VerifyEmailCode(req.Email, req.Code, c)
	if !valid {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  utils.ErrorMessage(code),
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

	config.Log.Infof("user %s bind email %s", userID, req.Email)

	c.JSON(http.StatusOK, gin.H{
		"code": utils.SUCCESS,
		"msg":  "绑定成功",
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

func MergeAccount(c *gin.Context) {
	var req MergeAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": utils.ErrCodeInvalidParams,
			"msg":  utils.ErrorMessage(utils.ErrCodeInvalidParams),
		})
		return
	}

	// 当前小程序用户ID
	tempUserID := getUserIdFromContext(c)

	// 1. 验证邮箱验证码
	valid, code := models.VerifyEmailCode(req.Email, req.Code, c)
	if !valid {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  utils.ErrorMessage(code),
		})
		return
	}

	// 2. 获取邮箱账号
	emailUser := models.QueryUser(req.Email)
	if emailUser == nil {
		c.JSON(http.StatusOK, gin.H{
			"code": utils.ErrCodeUserNotFound,
			"msg":  utils.ErrorMessage(utils.ErrCodeUserNotFound),
		})
		return
	}

	// 3. 检查邮箱账号是否已绑定其他微信
	if emailUser.WechatOpenID != nil && *emailUser.WechatOpenID != "" {
		c.JSON(http.StatusOK, gin.H{
			"code": utils.ErrCodeWechatAlreadyBound,
			"msg":  utils.ErrorMessage(utils.ErrCodeWechatAlreadyBound),
		})
		return
	}

	// 4. 获取小程序临时账号
	tempUser, getCode := models.GetUserById(tempUserID)
	if getCode != utils.SUCCESS {
		c.JSON(http.StatusOK, gin.H{
			"code": getCode,
			"msg":  utils.ErrorMessage(getCode),
		})
		return
	}

	// 5. 合并账号

	// 5.3 删除临时的小程序账号
	deleteCode := models.DeleteUser(tempUserID)
	if deleteCode != utils.SUCCESS {
		config.Log.Warnf("delete temp user failed: %s", tempUserID)
	}

	// 5.1 把openid绑定到邮箱账号
	if tempUser.WechatOpenID != nil {
		emailUser.WechatOpenID = tempUser.WechatOpenID
	}

	// 5.2 VIP状态合并（取最长的）
	if tempUser.VipExpiry.After(emailUser.VipExpiry) {
		emailUser.VipExpiry = tempUser.VipExpiry
	}

	// 5.4 更新邮箱账号
	err := config.DB.Save(emailUser).Error
	if err != nil {
		config.Log.Errorf("merge account save error: %v", err)
		c.JSON(http.StatusOK, gin.H{
			"code": utils.ErrCodeAccountMergeFailed,
			"msg":  utils.ErrorMessage(utils.ErrCodeAccountMergeFailed),
		})
		return
	}

	config.Log.Infof("account merged: temp user %s -> email user %s", tempUserID, emailUser.ID.String())

	// 6. Generate token version and store to Redis (按客户端类型区分，小程序端)
	tokenVersion := models.GenerateTokenVersion()
	if err := models.SetTokenVersion(c.Request.Context(), emailUser.ID, models.ClientTypeMiniprogram, tokenVersion); err != nil {
		config.Log.Errorf("set token version error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": utils.ErrCodeTokenGenerate,
			"msg":  "登录失败，请重试",
		})
		return
	}

	// 7. 生成新token（使用邮箱账号）
	token, err := utils.GenerateJWT(emailUser.ID, emailUser.UserName, emailUser.Role, emailUser.VipExpiry, emailUser.Email, tokenVersion, models.ClientTypeMiniprogram)
	if err != nil {
		config.Log.Errorf("generate token error: %v", err)
		c.JSON(http.StatusOK, gin.H{
			"code": utils.ErrCodeTokenGenerate,
			"msg":  utils.ErrorMessage(utils.ErrCodeTokenGenerate),
		})
		return
	}

	// 7. 返回合并后的用户信息
	c.JSON(http.StatusOK, gin.H{
		"code":  utils.SUCCESS,
		"token": token,
		"data": gin.H{
			"id":         emailUser.ID.String(),
			"username":   emailUser.UserName,
			"email":      emailUser.Email,
			"role":       emailUser.Role,
			"vip_expiry": emailUser.VipExpiry,
			"has_email":  true,
			"has_wechat": emailUser.WechatOpenID != nil && *emailUser.WechatOpenID != "",
		},
		"msg": "Account merged successfully",
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
