package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"uu/models"
	"uu/utils"
)

func UpdateSetting(c *gin.Context) {
	userId := getUserIdFromContext(c)
	var s models.Settings
	if err := c.ShouldBindJSON(&s); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": utils.InvalidParameter,
			"msg":  utils.ErrorMessage(utils.InvalidParameter),
		})
		return
	}
	code := models.UpdateSetting(userId, s)
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  utils.ErrorMessage(code),
	})
}

func GetSettings(c *gin.Context) {
	userId := getUserIdFromContext(c)
	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": utils.InvalidParameter,
			"msg":  utils.ErrorMessage(utils.InvalidParameter),
		})
		return
	}
	setting, code := models.GetUserSetting(userId)
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"data": setting,
		"msg":  utils.ErrorMessage(code),
	})
}

func getUserIdFromContext(c *gin.Context) string {
	val, exists := c.Get("userID")
	if !exists {
		return ""
	}

	switch v := val.(type) {
	case string:
		return v
	default:
		return ""
	}
}
