package api

import (
	"net/http"
	"strconv"
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
	var user = models.User{
		ID:       uid,
		UserName: reg.Username,
		Email:    reg.Email,
		Password: models.ScryptPw(reg.Password),
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
	c.JSON(http.StatusOK, gin.H{
		"code":     code,
		"username": user.UserName,
		"email":    user.Email,
		"msg":      utils.ErrorMessage(code),
	})
}

func GetUserList(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	pageNum, _ := strconv.Atoi(c.DefaultQuery("page_num", "1"))
	search := c.Query("search")
	users, total, code := models.GetUserList(pageSize, pageNum, search)
	c.JSON(http.StatusOK, gin.H{
		"code":      code,
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
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  utils.ErrorMessage(code),
	})
}

func GetSelfInfo(c *gin.Context) {
	userId, _ := c.Get("userID")
	username, _ := c.Get("username")
	role, _ := c.Get("role")
	vipExpiry, _ := c.Get("vipExpiry")
	email, _ := c.Get("email")
	c.JSON(http.StatusOK, gin.H{
		"code": 1,
		"msg":  "success",
		"data": gin.H{
			"id":         userId,
			"username":   username,
			"email":      email,
			"role":       role,
			"vip_expiry": vipExpiry,
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
	newExpiry, code := models.RenewVIP(req.UserId, req.Days)
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
	userIdStr := userID.(uuid.UUID).String()
	user, code := models.GetUserById(userIdStr)
	if code != utils.SUCCESS {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": code,
			"msg":  utils.ErrorMessage(code),
		})
		return
	}

	// 生成新的token（使用最新的用户信息）
	newToken, err := utils.GenerateJWT(user.ID, user.UserName, user.Role, user.VipExpiry, user.Email)
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
