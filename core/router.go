package core

import (
	"time"
	"uu/api"
	"uu/config"
	"uu/middleware"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	gin.SetMode(config.CONFIG.Server.Env)
	r := gin.Default()
	r.Use(gin.Recovery(), middleware.Cors(), middleware.Logger())

	// Gzip 压缩（节省带宽，压缩 JSON 响应）
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	v1 := r.Group("api/v1")

	v1.GET("captcha", middleware.RateLimiterByIP(middleware.RateLimiterConfig{
		Window:      60 * time.Second,
		MaxRequests: 10,
		KeyPrefix:   "user:captcha",
	}), api.GenerateCaptcha)

	// 公开接口
	public := v1.Group("public")
	{
		// 首页数据：每个IP每分钟最多30次
		public.GET("home",
			middleware.RateLimiterByIP(middleware.RateLimiterConfig{
				Window:      60 * time.Second,
				MaxRequests: 30,
				KeyPrefix:   "public:home",
			}),
			api.GetPublicHomeData,
		)
	}

	user := v1.Group("user")
	{
		user.POST("register",
			middleware.RateLimiterByIP(middleware.RateLimiterConfig{
				Window:      60 * time.Second,
				MaxRequests: 5,
				KeyPrefix:   "user:register",
			}),
			api.Register,
		)

		user.POST("login",
			middleware.RateLimiterByIP(middleware.RateLimiterConfig{
				Window:      60 * time.Second,
				MaxRequests: 10,
				KeyPrefix:   "user:login",
			}),
			api.Login,
		)
		user.POST("email-login",
			middleware.RateLimiterByIP(middleware.RateLimiterConfig{
				Window:      60 * time.Second,
				MaxRequests: 10,
				KeyPrefix:   "user:email-login",
			}),
			api.LoginByEmail,
		)

		user.POST("send-email",
			middleware.RateLimiterByIP(middleware.RateLimiterConfig{
				Window:      60 * time.Second,
				MaxRequests: 3,
				KeyPrefix:   "user:send-email",
			}),
			api.SendEmailCode,
		)

		user.POST("reset-password",
			middleware.RateLimiterByIP(middleware.RateLimiterConfig{
				Window:      3600 * time.Second,
				MaxRequests: 5,
				KeyPrefix:   "user:reset-password",
			}),
			api.ResetUserPassword,
		)

		user.POST("verify-email-code",
			middleware.RateLimiterByIP(middleware.RateLimiterConfig{
				Window:      60 * time.Second,
				MaxRequests: 10,
				KeyPrefix:   "user:verify-email-code",
			}),
			api.VerifyEmailCode,
		)
		user.POST("email-exist", api.JudgeEmail)
	}
	authUser := user.Group("")
	authUser.Use(middleware.AuthMiddleware())
	{
		authUser.GET("self", api.GetSelfInfo)
		authUser.PUT("name", api.UpdateUserName)
		authUser.POST("refresh", api.RefreshToken)
		authUser.POST("logout", api.Logout)
	}

	authGoods := v1.Group("goods")
	authGoods.Use(middleware.AuthMiddleware())
	authGoods.Use(middleware.RateLimiterByUser(middleware.RateLimiterConfig{
		Window:      60 * time.Second,
		MaxRequests: 30,
		KeyPrefix:   "auth:goods",
	}))
	{
		authGoods.GET("search", api.SearchGoods)            // 搜索饰品
		authGoods.GET("detail", api.GetGoodsDetail)         // 饰品详情
		authGoods.GET("related-wears", api.GetRelatedWears) // 同款不同磨损
	}

	vip := v1.Group("vip")
	vip.Use(middleware.AuthMiddleware(), middleware.AuthVIPMiddleware())
	goods := vip.Group("goods")

	goods.Use(middleware.RateLimiterByUser(middleware.RateLimiterConfig{
		Window:      60 * time.Second,
		MaxRequests: 60,
		KeyPrefix:   "vip:goods",
	}))
	{
		goods.GET("data", api.GetGoods)
		goods.GET("category", api.GetGoodsCategory)
		goods.GET("price-history", api.GetPriceHistory)
		goods.GET("price-increase", api.GetPriceIncreaseByU)
		goods.GET("big-item-bidding", api.GetBigItemBidding)
	}

	admin := v1.Group("admin")
	admin.Use(middleware.AuthMiddleware(), middleware.AuthAdminMiddleware())

	admin.Use(middleware.RateLimiterByUser(middleware.RateLimiterConfig{
		Window:      60 * time.Second,
		MaxRequests: 60,
		KeyPrefix:   "admin",
	}))
	{
		admin.GET("users", api.GetUserList)
		admin.DELETE("user", api.DeleteUser)
		admin.POST("vip-expiry", api.RenewVipExpiry)
		admin.GET("orders", api.GetAllOrders)
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

	// 支付相关API
	// 支付回调
	v1.POST("payment/notify", api.PayNotify)

	// 支付接口
	payment := v1.Group("payment")
	payment.Use(middleware.AuthMiddleware())
	{
		payment.GET("vip-price", api.GetVipPrice)         // 获取VIP价格
		payment.POST("create", api.CreatePayOrder)        // 创建支付订单
		payment.GET("query", api.QueryPayOrder)           // 查询订单状态
		payment.GET("orders", api.GetUserOrders)          // 获取订单列表
		payment.GET("vip-records", api.GetUserVipRecords) // 获取VIP开通记录
	}

	// 微信小程序相关API
	wechat := v1.Group("wechat")
	{
		// 微信登录：每个IP每分钟最多10次
		wechat.POST("login",
			middleware.RateLimiterByIP(middleware.RateLimiterConfig{
				Window:      60 * time.Second,
				MaxRequests: 10,
				KeyPrefix:   "wechat:login",
			}),
			api.WechatLogin,
		)
		// 绑定邮箱：每个用户每分钟最多5次
		wechat.POST("bind-email",
			middleware.AuthMiddleware(),
			middleware.RateLimiterByUser(middleware.RateLimiterConfig{
				Window:      60 * time.Second,
				MaxRequests: 5,
				KeyPrefix:   "wechat:bind-email",
			}),
			api.BindEmail,
		)
		// 合并账号：每个用户每小时最多3次
		wechat.POST("merge-account",
			middleware.AuthMiddleware(),
			middleware.RateLimiterByUser(middleware.RateLimiterConfig{
				Window:      3600 * time.Second,
				MaxRequests: 3,
				KeyPrefix:   "wechat:merge-account",
			}),
			api.MergeAccount,
		)
		// Web用户绑定微信：每个用户每分钟最多5次
		wechat.POST("bind-wechat",
			middleware.AuthMiddleware(),
			middleware.RateLimiterByUser(middleware.RateLimiterConfig{
				Window:      60 * time.Second,
				MaxRequests: 5,
				KeyPrefix:   "wechat:bind-wechat",
			}),
			api.BindWechat,
		)
		// 发送验证码：每个用户每分钟最多3次
		wechat.POST("send-email-code",
			middleware.AuthMiddleware(),
			middleware.RateLimiterByUser(middleware.RateLimiterConfig{
				Window:      60 * time.Second,
				MaxRequests: 3,
				KeyPrefix:   "wechat:send-email-code",
			}),
			api.SendEmailCode,
		)
	}

	return r
}
