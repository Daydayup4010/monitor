package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"uu/models"
	"uu/utils"
)

// LoginRequest 登录请求体
type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
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
	token, err := utils.GenerateJWT(user.ID, user.UserName, user.Role, user.VipExpiry, user.Email)
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
		"user": gin.H{
			"id":       user.ID,
			"username": user.UserName,
			"email":    user.Email,
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

	token, err := utils.GenerateJWT(user.ID, user.UserName, user.Role, user.VipExpiry, user.Email)
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
		"user": gin.H{
			"id":       user.ID,
			"username": user.UserName,
			"email":    user.Email,
		},
		"msg": utils.ErrorMessage(utils.SUCCESS),
	})
}
