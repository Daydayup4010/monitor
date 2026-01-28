package api

import (
	"net/http"
	"uu/models"
	"uu/utils"

	"github.com/gin-gonic/gin"
)

// GetMinAppConfig 获取小程序配置（公开API，无需登录）
func GetMinAppConfig(c *gin.Context) {
	cfg := models.GetMinAppConfig()
	c.JSON(http.StatusOK, gin.H{
		"code": utils.SUCCESS,
		"msg":  utils.ErrorMessage(utils.SUCCESS),
		"data": cfg,
	})
}

// SetMinAppVipEnabled 设置小程序VIP开通入口开关（管理员API）
func SetMinAppVipEnabled(c *gin.Context) {
	var req struct {
		Enabled bool `json:"enabled"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": utils.ErrCodeInvalidParams,
			"msg":  utils.ErrorMessage(utils.ErrCodeInvalidParams),
		})
		return
	}

	value := "0"
	if req.Enabled {
		value = "1"
	}

	err := models.SetSystemConfig(models.ConfigKeyMinAppVipEnabled, value, "")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": utils.ERROR,
			"msg":  "设置失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": utils.SUCCESS,
		"msg":  utils.ErrorMessage(utils.SUCCESS),
	})
}

// GetSystemConfigs 获取所有系统配置（管理员API）
func GetSystemConfigs(c *gin.Context) {
	cfg := models.GetMinAppConfig()
	c.JSON(http.StatusOK, gin.H{
		"code": utils.SUCCESS,
		"msg":  utils.ErrorMessage(utils.SUCCESS),
		"data": gin.H{
			"minapp_vip_enabled": cfg["vip_enabled"],
		},
	})
}
