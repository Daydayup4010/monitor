package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"uu/models"
)

func GetGoods(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.Query("page_size"))
	pageNum, _ := strconv.Atoi(c.Query("page_num"))
	userId := c.Query("user_id")
	sort := c.Query("sort")
	desc, _ := strconv.ParseBool(c.Query("desc"))
	category := c.Query("category")
	s, total, err := models.GetGoods(userId, pageSize, pageNum, desc, sort, category)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":  1,
		"data":  s,
		"total": total,
		"msg":   "success",
	})
}
