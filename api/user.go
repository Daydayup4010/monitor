package api

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strconv"
	"uu/config"
	"uu/models"
	"uu/utils"
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
		config.Log.Errorf("register user error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid parameter",
		})
		return
	}

	result, err := models.VerifyEmailCode(reg.Email, reg.Code, c.Request.Context())
	if err != nil || !result {
		config.Log.Errorf("The verification code is incorrect: %s", err)
		c.JSON(http.StatusOK, gin.H{
			"error": "verification code is incorrect",
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
	err = models.CreateUser(&user)
	if err != nil {
		config.Log.Errorf("Create User error: %v", err)
		c.JSON(http.StatusOK, gin.H{
			"error": "Create User error",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"username": user.UserName,
		"email":    user.Email,
		"msg":      "success",
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

func SendEmailCode(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": "invalid parameter",
		})
		return
	}
	code := utils.GenerateVerificationCode(6)
	err = models.SaveEmailCode(req.Email, code, c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "generate verify code fail",
		})
		return
	}
	err = config.CONFIG.Email.SendVerificationCode(req.Email, code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "send code fail",
		})
		config.Log.Errorf("send code fail: %v", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "success",
	})
}

func GetSelfInfo(c *gin.Context) {
	userId, _ := c.Get("userID")
	username, _ := c.Get("username")
	c.JSON(http.StatusOK, gin.H{
		"data": map[string]interface{}{
			"id":        userId,
			"user_name": username,
		},
	})
}

func UpdateUserName(c *gin.Context) {
	userId := c.Query("id")
	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": "invalid parameter",
		})
		return
	}
	err := models.UpdateUserName(req.Name, userId)
	if err != nil {
		config.Log.Errorf("update user name error :%v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "update user name fail",
		})
	}
}

func ResetUserPassword(c *gin.Context) {
	var req ResetPassword
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid parameter",
		})
		return
	}
	result, err := models.VerifyEmailCode(req.Email, req.Code, c.Request.Context())
	if err != nil || !result {
		c.JSON(http.StatusOK, gin.H{
			"error": "verification code is incorrect",
		})
		return
	}
	password := models.ScryptPw(req.Password)
	err = models.ResetPassword(req.Email, password)
	if err != nil {
		config.Log.Errorf("reset password fail: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "reset password fail",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "success",
	})

}
