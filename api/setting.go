package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"uu/config"
	"uu/models"
)

func UpdateSetting(c *gin.Context) {
	var s models.Settings
	_ = c.ShouldBindJSON(&s)
	err := s.UpdateSetting(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "update setting err",
		})
		config.Log.Errorf("update setting err: %s", err)
	} else {
		models.UpdateSkinItems()
		config.Log.Info("Update skin item")
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "success",
		})
	}

}

func GetSettings(c *gin.Context) {
	var s models.Settings
	err := s.GetSettings(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "get settings err",
		})
		config.Log.Errorf("get settings err: %s", err)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"data": s,
			"msg":  "success",
		})
	}
}
