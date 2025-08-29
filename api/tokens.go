package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"uu/config"
	"uu/models"
)

func UpdateUUToken(c *gin.Context) {
	var yp models.UUToken
	_ = c.ShouldBindJSON(&yp)
	err := yp.SetUUToken(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "update youpin token error",
		})
		config.Log.Errorf("update youpin token error: %s", err)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "success",
		})
	}
}

func GetUUToken(c *gin.Context) {
	var yp models.UUToken
	err := yp.GetUUToken(c)
	message := "success"
	if err != nil {
		config.Log.Errorf("get youpin token error: %s", err)
		message = "get youpin token error"
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    1,
		"data":    yp,
		"message": message,
	})
}
