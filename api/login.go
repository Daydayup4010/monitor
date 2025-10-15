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
			"error": err.Error(),
		})
		return
	}
	if !models.IfExistEmail(req.Email) {
		c.JSON(http.StatusOK, gin.H{
			"error": "user not exist",
		})
		return
	}
	user := models.QueryUser(req.Email)
	if models.ScryptPw(req.Password) != user.Password {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid Credential",
		})
		return
	}
	user.LastLogin = time.Now()
	models.UpdateUserLastLogin(user)
	token, err := utils.GenerateJWT(user.ID, user.UserName, user.Role)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": "generate token fail",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.UserName,
			"email":    user.Email,
		},
	})
}
