package core

import (
	"uu/api"
	"uu/config"
	"uu/middleware"

	"github.com/gin-gonic/gin"
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
		user.POST("email-login", api.LoginByEmail)
		user.POST("send-email", api.SendEmailCode)
		user.POST("reset-password", api.ResetUserPassword)
		user.POST("email-exist", api.JudgeEmail)
	}
	authUser := user.Group("")
	authUser.Use(middleware.AuthMiddleware())
	{
		authUser.GET("self", api.GetSelfInfo)
		authUser.PUT("name", api.UpdateUserName)
		authUser.POST("refresh", api.RefreshToken)
	}

	vip := v1.Group("vip")
	vip.Use(middleware.AuthMiddleware(), middleware.AuthVIPMiddleware())
	goods := vip.Group("goods")
	{
		goods.GET("data", api.GetGoods)
		goods.GET("category", api.GetGoodsCategory)
	}

	admin := v1.Group("admin")
	admin.Use(middleware.AuthMiddleware(), middleware.AuthAdminMiddleware())
	{
		admin.GET("users", api.GetUserList)
		admin.DELETE("user", api.DeleteUser)
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

	// 微信小程序相关API
	wechat := v1.Group("wechat")
	{
		wechat.POST("login", api.WechatLogin)                                          // 微信登录
		wechat.POST("bind-email", middleware.AuthMiddleware(), api.BindEmail)          // 绑定邮箱
		wechat.POST("merge-account", middleware.AuthMiddleware(), api.MergeAccount)    // 合并账号
		wechat.POST("bind-wechat", middleware.AuthMiddleware(), api.BindWechat)        // Web用户绑定微信
		wechat.POST("send-email-code", middleware.AuthMiddleware(), api.SendEmailCode) // 发送验证码
	}

	return r
}
