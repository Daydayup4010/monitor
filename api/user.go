package api

import (
	"net/http"
	"strconv"
	"time"
	"uu/config"
	"uu/models"
	"uu/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RegisterRequest 注册请求体
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Code     string `json:"code" binding:"required"`
}

type ResetPassword struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Code     string `json:"code" binding:"required"`
}

func Register(c *gin.Context) {
	var reg RegisterRequest
	if err := c.ShouldBindJSON(&reg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": utils.InvalidParameter,
			"msg":  utils.ErrorMessage(utils.InvalidParameter),
		})
		return
	}

	result, code := models.VerifyEmailCode(reg.Email, reg.Code, c.Request.Context())
	if !result {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  utils.ErrorMessage(code),
		})
		return
	}

	if models.IfExistEmail(reg.Email) {
		c.JSON(http.StatusOK, gin.H{
			"code": utils.ErrCodeEmailTaken,
			"msg":  utils.ErrorMessage(utils.ErrCodeEmailTaken),
		})
		return
	}
	uid, _ := uuid.NewV7()
	newExpiry := time.Now().AddDate(0, 0, 2) // 新用户免费3天VIP
	var user = models.User{
		ID:        uid,
		UserName:  reg.Username,
		Email:     reg.Email,
		Password:  models.ScryptPw(reg.Password),
		Role:      1,
		VipExpiry: newExpiry,
	}

	code = models.CreateDefaultSetting(uid.String())
	if code != utils.SUCCESS {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": code,
			"msg":  utils.ErrorMessage(code),
		})
		return
	}
	code = models.CreateUser(&user)
	if code != utils.SUCCESS {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  utils.ErrorMessage(code),
		})
		return
	}

	// 注册成功后自动登录：生成 token（Web端注册）
	tokenVersion := models.GenerateTokenVersion()
	if err := models.SetTokenVersion(c.Request.Context(), user.ID, models.ClientTypeWeb, tokenVersion); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": utils.SUCCESS,
			"msg":  "注册成功，请登录",
		})
		return
	}

	token, err := utils.GenerateJWT(user.ID, user.UserName, user.Role, user.VipExpiry, user.Email, tokenVersion, models.ClientTypeWeb)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": utils.SUCCESS,
			"msg":  "注册成功，请登录",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":  utils.SUCCESS,
		"token": token,
		"data": gin.H{
			"id":         user.ID,
			"username":   user.UserName,
			"email":      user.Email,
			"role":       user.Role,
			"vip_expiry": user.VipExpiry,
		},
		"msg": utils.ErrorMessage(utils.SUCCESS),
	})
}

func GetUserList(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	pageNum, _ := strconv.Atoi(c.DefaultQuery("page_num", "1"))
	search := c.Query("search")
	users, total, vipCount, code := models.GetUserList(pageSize, pageNum, search)
	c.JSON(http.StatusOK, gin.H{
		"code":      code,
		"data":      users,
		"total":     total,
		"vip_count": vipCount,
		"page_size": pageSize,
		"page_num":  pageNum,
	})
}

func SendEmailCode(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": utils.InvalidParameter,
			"msg":  utils.ErrorMessage(utils.InvalidParameter),
		})
		return
	}
	verifyCode := utils.GenerateVerificationCode(6)
	code := models.SaveEmailCode(req.Email, verifyCode, c.Request.Context())
	if code != utils.SUCCESS {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": code,
			"msg":  utils.ErrorMessage(code),
		})
		return
	}
	code = config.CONFIG.Email.SendVerificationCode(req.Email, verifyCode)

	// 检查邮箱是否已存在（用于小程序绑定邮箱时判断是新绑定还是合并账号）
	emailExists := models.IfExistEmail(req.Email)

	c.JSON(http.StatusOK, gin.H{
		"code":         code,
		"msg":          utils.ErrorMessage(code),
		"email_exists": emailExists,
	})
}

// VerifyEmailCode 验证邮箱验证码（不删除验证码，用于分步验证）
func VerifyEmailCode(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
		Code  string `json:"code" binding:"required,len=6"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": utils.InvalidParameter,
			"msg":  utils.ErrorMessage(utils.InvalidParameter),
		})
		return
	}

	valid, code := models.CheckEmailCode(req.Email, req.Code, c.Request.Context())
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": code,
			"msg":  utils.ErrorMessage(code),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 1,
		"msg":  "验证成功",
	})
}

func GetSelfInfo(c *gin.Context) {
	userId, _ := c.Get("userID")

	// 从数据库获取最新的用户信息（而不是从 token 读取，确保 VIP 状态是最新的）
	user, code := models.GetUserById(userId.(uuid.UUID).String())
	if code != utils.SUCCESS {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  "获取用户信息失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": utils.SUCCESS,
		"msg":  "success",
		"data": gin.H{
			"id":         user.ID,
			"username":   user.UserName,
			"email":      user.Email,
			"role":       user.Role,
			"vip_expiry": user.VipExpiry,
		},
	})
}

func UpdateUserName(c *gin.Context) {
	userId := c.Query("id")
	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": utils.InvalidParameter,
			"msg":  utils.ErrorMessage(utils.InvalidParameter),
		})
		return
	}
	code := models.UpdateUserName(req.Name, userId)
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  utils.ErrorMessage(code),
	})

}

func ResetUserPassword(c *gin.Context) {
	var req ResetPassword
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": utils.InvalidParameter,
			"msg":  utils.ErrorMessage(utils.InvalidParameter),
		})
		return
	}
	result, code := models.VerifyEmailCode(req.Email, req.Code, c.Request.Context())
	if !result {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  utils.ErrorMessage(code),
		})
		return
	}
	password := models.ScryptPw(req.Password)
	code = models.ResetPassword(req.Email, password)
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  utils.ErrorMessage(code),
	})

}

func DeleteUser(c *gin.Context) {
	userId := c.Query("user_id")
	code := models.DeleteUser(userId)
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  utils.ErrorMessage(code),
	})
}

func RenewVipExpiry(c *gin.Context) {
	var req struct {
		UserId string `json:"user_id"`
		Days   int    `json:"days"` // 注意：字段名为days，但实际表示月数
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": utils.InvalidParameter,
			"msg":  utils.ErrorMessage(utils.InvalidParameter),
		})
		return
	}
	// Days参数实际表示月数
	newExpiry, userEmail, code := models.RenewVIP(req.UserId, req.Days)

	// 如果续费成功，发送邮件通知用户
	if code == utils.SUCCESS && userEmail != "" {
		go func() {
			expiryStr := newExpiry.Format("2006-01-02 15:04:05")
			emailCode := config.CONFIG.Email.SendVIPNotification(userEmail, req.Days, expiryStr)
			if emailCode != utils.SUCCESS {
				config.Log.Errorf("send vip notification email failed for user: %s", userEmail)
			}
		}()
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  "success",
		"date": newExpiry,
	})
}

// RefreshToken 刷新Token
func RefreshToken(c *gin.Context) {
	// 从context获取当前用户信息（已通过AuthMiddleware验证）
	userID, _ := c.Get("userID")
	username, _ := c.Get("username")
	email, _ := c.Get("email")

	// 从数据库获取最新的用户信息（确保role和vipExpiry是最新的）
	userIDVal, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": utils.ErrCodeInvalidToken,
			"msg":  utils.ErrorMessage(utils.ErrCodeInvalidToken),
		})
		return
	}
	userIdStr := userIDVal.String()
	user, code := models.GetUserById(userIdStr)
	if code != utils.SUCCESS {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": code,
			"msg":  utils.ErrorMessage(code),
		})
		return
	}

	// 获取客户端类型
	clientType, _ := c.Get("clientType")
	clientTypeStr, ok := clientType.(string)
	if !ok || clientTypeStr == "" {
		clientTypeStr = models.ClientTypeWeb
	}

	// 获取当前的 token 版本号（沿用，不重新生成）
	currentVersion, err := models.GetTokenVersion(c.Request.Context(), user.ID, clientTypeStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": utils.ErrCodeTokenKicked,
			"msg":  "登录已过期，请重新登录",
		})
		return
	}

	// 刷新 Redis 中版本号的过期时间
	if err := models.SetTokenVersion(c.Request.Context(), user.ID, clientTypeStr, currentVersion); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": utils.ErrCodeTokenGenerate,
			"msg":  "刷新失败，请重试",
		})
		return
	}

	// 生成新的token（使用最新的用户信息 + 原版本号）
	newToken, err := utils.GenerateJWT(user.ID, user.UserName, user.Role, user.VipExpiry, user.Email, currentVersion, clientTypeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": utils.ErrCodeTokenGenerate,
			"msg":  utils.ErrorMessage(utils.ErrCodeTokenGenerate),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":  utils.SUCCESS,
		"msg":   utils.ErrorMessage(utils.SUCCESS),
		"token": newToken,
		"data": gin.H{
			"id":         user.ID,
			"username":   username,
			"email":      email,
			"role":       user.Role,
			"vip_expiry": user.VipExpiry,
		},
	})
}

// Logout 登出（使当前客户端的 token 失效）
func Logout(c *gin.Context) {
	userID, _ := c.Get("userID")
	userIDVal, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": utils.ErrCodeInvalidToken,
			"msg":  utils.ErrorMessage(utils.ErrCodeInvalidToken),
		})
		return
	}

	// 获取客户端类型
	clientType, _ := c.Get("clientType")
	clientTypeStr, ok := clientType.(string)
	if !ok || clientTypeStr == "" {
		clientTypeStr = models.ClientTypeWeb
	}

	// 删除 Redis 中当前客户端类型的 token 版本号
	if err := models.InvalidateTokenVersion(c.Request.Context(), userIDVal, clientTypeStr); err != nil {
		config.Log.Errorf("logout invalidate token error: %v", err)
		// 即使失败也返回成功，因为客户端会清除本地 token
	}

	c.JSON(http.StatusOK, gin.H{
		"code": utils.SUCCESS,
		"msg":  "登出成功",
	})
}

func JudgeEmail(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": utils.InvalidParameter,
			"msg":  utils.ErrorMessage(utils.InvalidParameter),
		})
		return
	}
	exist := models.IfExistEmail(req.Email)
	if exist {
		c.JSON(http.StatusOK, gin.H{
			"code": utils.ErrCodeEmailTaken,
			"msg":  utils.ErrorMessage(utils.ErrCodeEmailTaken),
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": utils.SUCCESS,
			"msg":  utils.ErrorMessage(utils.SUCCESS),
		})
		return
	}

}
