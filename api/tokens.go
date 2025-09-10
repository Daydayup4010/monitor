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

func UpdateBuffToken(c *gin.Context) {
	var buff models.BuffToken
	_ = c.ShouldBindJSON(&buff)
	err := buff.SetBuffToken(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "update buff token error",
		})
		config.Log.Errorf("update buff token error: %s", err)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "success",
		})
	}

}

func GetVerify(c *gin.Context) {
	var uu models.UUToken
	var buff models.BuffToken
	var expired = map[string]string{
		"uu":   "yes",
		"buff": "yes",
	}
	err := uu.GetUUExpired()
	err = buff.GetBuffExpired()
	if err != nil {
		config.Log.Errorf("Get token expired error : %v", err)
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": expired,
			"msg":  "Get token expired error",
		})
	} else {
		expired["uu"] = uu.Expired
		expired["buff"] = buff.Expired
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"data": expired,
			"msg":  "success",
		})
	}
}

func VerifyToken(c *gin.Context) {
	var uu models.UUToken
	var buff models.BuffToken
	err := uu.UpdateUUExpired()
	err = buff.UpdateBuffExpired()
	if err != nil {
		config.Log.Errorf("Update token expired error : %v", err)
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "Update token expired error",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "success",
		})
	}
}
