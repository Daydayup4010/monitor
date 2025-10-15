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
	r.POST("api/register", api.Register)
	r.POST("api/login", api.Login)
	v1 := r.Group("api/v1")
	{

	}
	v1.Use(middleware.AuthMiddleware())
	vip := v1.Group("vip")
	vip.Use(middleware.AuthVIPMiddleware())
	{
		vip.GET("data", api.GetSkinItem)
	}
	admin := v1.Group("admin")
	admin.Use(middleware.AuthAdminMiddleware())
	tokens := admin.Group("tokens")
	{
		tokens.POST("uu", api.UpdateUUToken)
		tokens.POST("buff", api.UpdateBuffToken)
		tokens.POST("verify", api.VerifyToken)
		tokens.GET("verify", api.GetVerify)

	}
	settings := admin.Group("settings")
	{
		settings.GET("", api.GetSettings)
		settings.PUT("", api.UpdateSetting)
		settings.POST("full_update", api.UpdateFull)
	}
	users := admin.Group("users")
	{
		users.GET("", api.GetUserList)
	}
	return r
}
