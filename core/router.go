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

	user := v1.Group("user")
	{
		user.POST("register", api.Register)
		user.POST("login", api.Login)
		user.POST("send-email", api.SendEmailCode)
		user.POST("reset-password", api.ResetUserPassword)
	}

	v1.Use(middleware.AuthMiddleware())
	{
		user.GET("self", api.GetSelfInfo)
		user.PUT("name", api.UpdateUserName)
	}

	vip := v1.Group("vip")
	vip.Use(middleware.AuthVIPMiddleware())
	{
		vip.GET("data", api.GetGoods)
	}

	admin := v1.Group("admin")
	admin.Use(middleware.AuthAdminMiddleware())
	{
		admin.GET("users", api.GetUserList)
		admin.DELETE("user", api.DeleteUser)
		admin.POST("full-update", api.UpdateFull)
		admin.POST("vip-expiry", api.RenewVipExpiry)
	}

	tokens := admin.Group("tokens")
	{
		tokens.POST("uu", api.UpdateUUToken)
		tokens.POST("buff", api.UpdateBuffToken)
		tokens.POST("verify", api.VerifyToken)
		tokens.GET("verify", api.GetVerify)

	}

	settings := vip.Group("settings")
	{
		settings.GET("", api.GetSettings)
		settings.PUT("", api.UpdateSetting)
	}
	return r
}
