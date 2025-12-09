package api

import (
	"context"
	"net/http"
	"time"
	"uu/config"
	"uu/utils"

	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
)

// Redis 存储验证码
type RedisCaptchaStore struct{}

// Set 存储验证码到 Redis
func (r *RedisCaptchaStore) Set(id string, value string) error {
	ctx := context.Background()
	// 验证码5分钟过期
	return config.RDB.Set(ctx, "captcha:"+id, value, 5*time.Minute).Err()
}

// Get 从 Redis 获取验证码
func (r *RedisCaptchaStore) Get(id string, clear bool) string {
	ctx := context.Background()
	key := "captcha:" + id
	val, err := config.RDB.Get(ctx, key).Result()
	if err != nil {
		return ""
	}
	if clear {
		config.RDB.Del(ctx, key) // 验证后删除，防止重复使用
	}
	return val
}

// Verify 验证验证码
func (r *RedisCaptchaStore) Verify(id, answer string, clear bool) bool {
	stored := r.Get(id, clear)
	if stored == "" {
		return false
	}
	return stored == answer
}

// 全局验证码存储
var CaptchaStore base64Captcha.Store = &RedisCaptchaStore{}

// GenerateCaptcha 生成图形验证码
func GenerateCaptcha(c *gin.Context) {
	// 配置验证码：数字，4位
	driver := base64Captcha.NewDriverDigit(
		80,  // 高度
		200, // 宽度
		4,   // 验证码长度
		0.7, // 干扰线数量
		80,  // 背景圆点数量
	)

	captcha := base64Captcha.NewCaptcha(driver, CaptchaStore)
	id, b64s, _, err := captcha.Generate()

	if err != nil {
		config.Log.Errorf("Generate captcha error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": utils.ErrCodeCaptchaGenerate,
			"msg":  utils.ErrorMessage(utils.ErrCodeCaptchaGenerate),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":        utils.SUCCESS,
		"captcha_id":  id,
		"captcha_img": b64s,
		"msg":         utils.ErrorMessage(utils.SUCCESS),
	})
}

// VerifyCaptcha 验证验证码（内部使用）
func VerifyCaptcha(captchaId, captchaCode string) bool {
	if captchaId == "" || captchaCode == "" {
		return false
	}
	return CaptchaStore.Verify(captchaId, captchaCode, true)
}
