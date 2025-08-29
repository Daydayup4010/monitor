package api

import (
	"github.com/gin-gonic/gin"
	"uu/services"
)

func GetItems(c *gin.Context) {
	items := services.GetUUItems(20, 1)
	c.JSON(200, gin.H{
		"data": items,
	})
}
