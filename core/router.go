package core

import (
	"github.com/gin-gonic/gin"
	"uu/api"
	"uu/config"
	"uu/middleware"
)

func InitRouter() *gin.Engine {
	gin.SetMode(config.CONFIG.Server.Env)
	r := gin.Default()
	r.Use(gin.Recovery(), middleware.Cors(), middleware.Logger())
	v1 := r.Group("api/v1")
	tokens := v1.Group("tokens")
	{
		tokens.POST("uu", api.UpdateUUToken)
		tokens.POST("buff", api.UpdateBuffToken)
		tokens.GET("verify", api.VerifyToken)

	}
	settings := v1.Group("settings")
	{
		settings.GET("", api.GetSettings)
		settings.PUT("", api.UpdateSetting)
	}
	data := v1.Group("data")
	{
		data.GET("", api.GetSkinItem)
		data.POST("full_update", api.UpdateFull)
	}
	return r
}
