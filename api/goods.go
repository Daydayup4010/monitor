package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"uu/models"
	"uu/utils"
)

func GetGoods(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.Query("page_size"))
	pageNum, _ := strconv.Atoi(c.Query("page_num"))
	userId := getUserIdFromContext(c)
	sort := c.Query("sort")
	desc, _ := strconv.ParseBool(c.Query("desc"))
	category := c.Query("category")
	s, total, code := models.GetGoods(userId, pageSize, pageNum, desc, sort, category)
	c.JSON(http.StatusOK, gin.H{
		"code":  code,
		"data":  s,
		"total": total,
		"msg":   utils.ErrorMessage(code),
	})
}
