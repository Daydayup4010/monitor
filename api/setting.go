package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"uu/config"
	"uu/models"
)

func UpdateSetting(c *gin.Context) {
	userId := c.Query("user_id")
	var s models.Settings
	_ = c.ShouldBindJSON(&s)
	err := models.UpdateSetting(userId, s)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "update setting err",
		})
		config.Log.Errorf("update setting err: %s", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 1,
		"msg":  "success",
	})
}

func GetSettings(c *gin.Context) {
	userId := c.Query("user_id")
	if userId == "" {
		c.JSON(http.StatusOK, gin.H{
			"error": "invalid parameter",
		})
		return
	}
	setting, err := models.GetUserSetting(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 0,
			"msg":  "get settings err",
		})
		config.Log.Errorf("get settings err: %s", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 1,
		"data": setting,
		"msg":  "success",
	})
}
