package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"uu/config"
	"uu/utils"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// RateLimiterConfig 限流配置
type RateLimiterConfig struct {
	// 时间窗口
	Window time.Duration
	// 窗口内最大请求数
	MaxRequests int
	// 限流的键前缀（用于区分不同接口）
	KeyPrefix string
}

// RateLimiterByUser 按用户ID限流（需要先经过 AuthMiddleware）
func RateLimiterByUser(cfg RateLimiterConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			// 如果未认证，降级为IP限流
			userID = getRealIP(c)
		}

		key := fmt.Sprintf("rate_limit:%s:user:%v", cfg.KeyPrefix, userID)
		handleRateLimit(c, key, cfg)
	}
}

// RateLimiterByIP 按IP限流（适合公开接口）
func RateLimiterByIP(cfg RateLimiterConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := getRealIP(c)
		key := fmt.Sprintf("rate_limit:%s:ip:%s", cfg.KeyPrefix, ip)
		handleRateLimit(c, key, cfg)
	}
}

// getRealIP 获取客户端真实IP
func getRealIP(c *gin.Context) string {
	// 优先从 X-Forwarded-For 获取（如果使用了反向代理）
	if xff := c.GetHeader("X-Forwarded-For"); xff != "" {
		// X-Forwarded-For 可能包含多个IP，取第一个
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// 其次从 X-Real-IP 获取
	if xri := c.GetHeader("X-Real-IP"); xri != "" {
		return xri
	}

	// 最后使用 RemoteAddr
	return c.ClientIP()
}

// handleRateLimit 统一的限流处理逻辑
func handleRateLimit(c *gin.Context, key string, cfg RateLimiterConfig) {
	allowed, remaining, resetTime, err := checkRateLimit(key, cfg.Window, cfg.MaxRequests)
	if err != nil {
		config.Log.Errorf("Rate limit check error: %v", err)
		// 出错时允许通过
		c.Next()
		return
	}

	// 设置响应头，告知客户端限流信息
	c.Header("X-RateLimit-Limit", strconv.Itoa(cfg.MaxRequests))
	c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))
	c.Header("X-RateLimit-Reset", strconv.FormatInt(resetTime, 10))

	if !allowed {
		c.JSON(http.StatusTooManyRequests, gin.H{
			"code":       utils.ErrCodeRateLimitExceeded,
			"msg":        utils.ErrorMessage(utils.ErrCodeRateLimitExceeded),
			"retryAfter": resetTime - time.Now().Unix(),
		})
		c.Abort()
		return
	}

	c.Next()
}

// checkRateLimit 使用滑动窗口算法检查限流
// 返回: (是否允许, 剩余次数, 重置时间戳, 错误)
func checkRateLimit(key string, window time.Duration, maxRequests int) (bool, int, int64, error) {
	ctx := context.Background()
	now := time.Now()
	windowStart := now.Add(-window).UnixMilli()
	nowMilli := now.UnixMilli()

	pipe := config.RDB.Pipeline()

	// 1. 删除窗口之外的旧记录
	pipe.ZRemRangeByScore(ctx, key, "0", strconv.FormatInt(windowStart, 10))

	// 2. 统计当前窗口内的请求数
	zcard := pipe.ZCard(ctx, key)

	// 3. 添加当前请求（使用时间戳作为score，纳秒时间戳作为member确保唯一性）
	member := fmt.Sprintf("%d:%d", nowMilli, now.UnixNano())
	pipe.ZAdd(ctx, key, redis.Z{Score: float64(nowMilli), Member: member})

	// 4. 设置key过期时间
	pipe.Expire(ctx, key, window+time.Second)

	// 执行pipeline
	_, err := pipe.Exec(ctx)
	if err != nil {
		return false, 0, 0, err
	}

	// 获取当前请求数（执行前的数量 + 1）
	count := int(zcard.Val()) + 1
	remaining := maxRequests - count
	if remaining < 0 {
		remaining = 0
	}

	// 计算重置时间
	resetTime := now.Add(window).Unix()

	// 判断是否超过限制
	allowed := count <= maxRequests

	return allowed, remaining, resetTime, nil
}
