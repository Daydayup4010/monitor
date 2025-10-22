package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"uu/services"
	"uu/utils"
)

func UpdateFull(c *gin.Context) {
	code := services.UpdateFullData()
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  utils.ErrorMessage(code),
	})
}
