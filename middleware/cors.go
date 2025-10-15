package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

func Cors() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins: []string{"*"},

		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},

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
