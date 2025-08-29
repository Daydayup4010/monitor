package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

func Cors() gin.HandlerFunc {
	return cors.New(cors.Config{
		// 生产环境建议明确指定允许的域名（不要用 *）
		AllowOrigins: []string{"*"},

		// 明确声明允许的方法
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},

		// 必须包含所有需要的自定义头
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
			"X-Requested-With",
		},

		// 暴露给前端可访问的响应头
		ExposeHeaders: []string{"Content-Length", "Authorization"},

		// 是否允许携带 Cookie（如果开启，AllowOrigins 不能为 *）
		AllowCredentials: false,

		// 预检请求缓存时间
		MaxAge: 12 * time.Hour,
	})
}
