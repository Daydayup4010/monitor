package api

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strconv"
	"strings"
	"uu/config"
	"uu/models"
)

// RegisterRequest 注册请求体
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func Register(c *gin.Context) {
	var reg RegisterRequest
	if err := c.ShouldBindJSON(&reg); err != nil {
		if strings.Contains(err.Error(), "RegisterRequest.Email") {

		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	if models.IfExistEmail(reg.Email) {
		c.JSON(http.StatusOK, gin.H{
			"error": "Email Registered",
		})
		return
	}
	uid, _ := uuid.NewV7()
	var user = models.User{
		ID:       uid,
		UserName: reg.Username,
		Email:    reg.Email,
		Password: models.ScryptPw(reg.Password),
	}
	err := models.CreateUser(&user)
	if err != nil {
		config.Log.Errorf("Create User error: %v", err)
		c.JSON(http.StatusOK, gin.H{
			"error": "Create User error",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":     "success",
		"username": user.UserName,
		"email":    user.Email,
	})
}

func GetUserList(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	pageNum, _ := strconv.Atoi(c.DefaultQuery("page_num", "1"))
	users, total := models.GetUserList(pageSize, pageNum)
	c.JSON(http.StatusOK, gin.H{
		"data":      users,
		"total":     total,
		"page_size": pageSize,
		"page_num":  pageNum,
	})
}
