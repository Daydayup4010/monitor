package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"uu/models"
	"uu/services"
	"uu/utils"
)

func UpdateUUToken(c *gin.Context) {
	var uu models.UUToken
	_ = c.ShouldBindJSON(&uu)
	code := uu.SetUUToken(c.Request.Context())
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  utils.ErrorMessage(code),
	})
}

func UpdateBuffToken(c *gin.Context) {
	var buff models.BuffToken
	_ = c.ShouldBindJSON(&buff)
	code := buff.SetBuffToken(c.Request.Context())
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  utils.ErrorMessage(code),
	})

}

func GetVerify(c *gin.Context) {
	var uu models.UUToken
	var buff models.BuffToken
	var expired = map[string]string{
		"uu":   "yes",
		"buff": "yes",
	}
	code := uu.GetUUExpired()
	code = buff.GetBuffExpired()
	expired["uu"] = uu.Expired
	expired["buff"] = buff.Expired
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"data": expired,
		"msg":  utils.ErrorMessage(code),
	})
}

func VerifyToken(c *gin.Context) {
	var uu models.UUToken
	var buff models.BuffToken
	var expired = map[string]string{
		"uu":   "yes",
		"buff": "yes",
	}
	services.VerifyUUToken()
	services.VerifyBuffToken()
	code := uu.GetUUExpired()
	code = buff.GetBuffExpired()
	expired["uu"] = uu.Expired
	expired["buff"] = buff.Expired
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"data": expired,
		"msg":  utils.ErrorMessage(code),
	})
}
