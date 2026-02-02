package api

import (
	"net/http"
	"time"
	"uu/models"
	"uu/utils"

	"github.com/gin-gonic/gin"
)

// LoginRequest 登录请求体
type LoginRequest struct {
	Email       string `json:"email" binding:"required"`
	Password    string `json:"password" binding:"required"`
	CaptchaId   string `json:"captcha_id" binding:"required"`
	CaptchaCode string `json:"captcha_code" binding:"required"`
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": utils.InvalidParameter,
			"msg":  utils.ErrorMessage(utils.InvalidParameter),
		})
		return
	}

	// 验证图形验证码
	if !VerifyCaptcha(req.CaptchaId, req.CaptchaCode) {
		c.JSON(http.StatusOK, gin.H{
			"code": utils.ErrCodeCaptchaInvalid,
			"msg":  utils.ErrorMessage(utils.ErrCodeCaptchaInvalid),
		})
		return
	}

	if !models.IfExistEmail(req.Email) {
		c.JSON(http.StatusOK, gin.H{
			"code": utils.ErrCodeUserNotFound,
			"msg":  utils.ErrorMessage(utils.ErrCodeUserNotFound),
		})
		return
	}
	user := models.QueryUser(req.Email)
	if models.ScryptPw(req.Password) != user.Password {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": utils.ErrCodeInvalidPassword,
			"msg":  utils.ErrorMessage(utils.ErrCodeInvalidPassword),
		})
		return
	}
	user.LastLogin = time.Now()
	models.UpdateUserLastLogin(user)

	// 生成 token 版本号并存储到 Redis（按客户端类型区分，Web端）
	tokenVersion := models.GenerateTokenVersion()
	if err := models.SetTokenVersion(c.Request.Context(), user.ID, models.ClientTypeWeb, tokenVersion); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": utils.ErrCodeTokenGenerate,
			"msg":  "登录失败，请重试",
		})
		return
	}

	token, err := utils.GenerateJWT(user.ID, user.UserName, user.Role, user.VipExpiry, user.Email, tokenVersion, models.ClientTypeWeb)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": utils.ErrCodeTokenGenerate,
			"msg":  utils.ErrorMessage(utils.ErrCodeTokenGenerate),
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

type EmailLoginRequest struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required"`
}

func LoginByEmail(c *gin.Context) {
	var req EmailLoginRequest
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

	if !models.IfExistEmail(req.Email) {
		c.JSON(http.StatusOK, gin.H{
			"code": utils.ErrCodeUserNotFound,
			"msg":  utils.ErrorMessage(utils.ErrCodeUserNotFound),
		})
		return
	}

	user := models.QueryUser(req.Email)
	user.LastLogin = time.Now()
	models.UpdateUserLastLogin(user)

	// 生成 token 版本号并存储到 Redis（按客户端类型区分，Web端）
	tokenVersion := models.GenerateTokenVersion()
	if err := models.SetTokenVersion(c.Request.Context(), user.ID, models.ClientTypeWeb, tokenVersion); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": utils.ErrCodeTokenGenerate,
			"msg":  "登录失败，请重试",
		})
		return
	}

	token, err := utils.GenerateJWT(user.ID, user.UserName, user.Role, user.VipExpiry, user.Email, tokenVersion, models.ClientTypeWeb)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": utils.ErrCodeTokenGenerate,
			"msg":  utils.ErrorMessage(utils.ErrCodeTokenGenerate),
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
