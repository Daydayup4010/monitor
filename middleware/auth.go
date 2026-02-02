package middleware

import (
	"net/http"
	"strings"
	"time"
	"uu/models"
	"uu/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":  utils.ErrCodeTokenNotFound,
				"error": "Authorization header required",
			})
			c.Abort()
			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":  utils.ErrCodeInvalidTokenFormat,
				"error": "Invalid authorization format",
			})
			c.Abort()
			return
		}
		claims, err := utils.ParseJWT(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":  utils.ErrCodeInvalidToken,
				"error": "invalid token",
			})
			c.Abort()
			return
		}

		// 获取客户端类型（兼容旧token，默认为web）
		clientType := claims.ClientType
		if clientType == "" {
			clientType = models.ClientTypeWeb
		}

		// 验证 token 版本号（按客户端类型检查，同类型客户端只能一个在线）
		if !models.ValidateTokenVersion(c.Request.Context(), claims.UserID, clientType, claims.TokenVersion) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":  utils.ErrCodeTokenKicked,
				"error": utils.ErrorMessage(utils.ErrCodeTokenKicked),
			})
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Set("vipExpiry", claims.VipExpiry)
		c.Set("email", claims.Email)
		c.Set("tokenVersion", claims.TokenVersion)
		c.Set("clientType", clientType)
		c.Next()
	}
}

func AuthVIPMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := getRoleFromContext(c)
		expiry := getExpiryFromContext(c)
		if !models.CanAccessVIPContent(role, expiry) {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "vip access required",
				"code":  "VIP_REQUIRED",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

func AuthAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := getRoleFromContext(c)
		if role != models.RoleAdmin {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Admin access required",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

func getRoleFromContext(c *gin.Context) int64 {
	val, exists := c.Get("role")
	if !exists {
		return 0
	}

	switch v := val.(type) {
	case int:
		return int64(v)
	case int64:
		return v
	case float64:
		return int64(v)
	default:
		return 0
	}
}

func getExpiryFromContext(c *gin.Context) time.Time {
	val, exists := c.Get("vipExpiry")
	if !exists {
		return time.Time{}
	}

	switch v := val.(type) {
	case time.Time:
		return v
	case string:
		t, _ := time.Parse(time.RFC3339, v)
		return t
	case int64:
		return time.Unix(v, 0)
	default:
		return time.Time{}
	}
}
