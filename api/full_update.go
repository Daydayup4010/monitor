package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"uu/services"
)

func UpdateFull(c *gin.Context) {
	services.UpdateFullData()
	c.JSON(http.StatusOK, gin.H{
		"msg": "success",
	})
}
