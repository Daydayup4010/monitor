package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"uu/models"
)

func GetSkinItem(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.Query("page_size"))
	pageNum, _ := strconv.Atoi(c.Query("page_num"))
	sort := c.Query("sort")
	desc, _ := strconv.ParseBool(c.Query("desc"))
	category := c.Query("category")
	s, total := models.GetSkinItems(pageSize, pageNum, desc, sort, category)
	c.JSON(http.StatusOK, gin.H{
		"code":  1,
		"data":  s,
		"total": total,
		"msg":   "success",
	})
}
