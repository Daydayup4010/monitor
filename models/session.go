package models

import (
	"context"
	"fmt"
	"time"
	"uu/config"

	"github.com/google/uuid"
)

const tokenVersionKeyPrefix = "user_token_version:"

// GenerateTokenVersion 生成新的 token 版本号（使用时间戳）
func GenerateTokenVersion() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// SetTokenVersion 存储用户的当前 token 版本号
func SetTokenVersion(ctx context.Context, userID uuid.UUID, version string) error {
	key := fmt.Sprintf("%s%s", tokenVersionKeyPrefix, userID.String())
	// 设置 25 小时过期（比 token 有效期稍长）
	return config.RDB.Set(ctx, key, version, 25*time.Hour).Err()
}

// GetTokenVersion 获取用户当前的 token 版本号
func GetTokenVersion(ctx context.Context, userID uuid.UUID) (string, error) {
	key := fmt.Sprintf("%s%s", tokenVersionKeyPrefix, userID.String())
	return config.RDB.Get(ctx, key).Result()
}

// ValidateTokenVersion 验证 token 版本号是否有效
func ValidateTokenVersion(ctx context.Context, userID uuid.UUID, version string) bool {
	currentVersion, err := GetTokenVersion(ctx, userID)
	if err != nil {
		return false
	}
	return currentVersion == version
}

// InvalidateTokenVersion 使 token 版本失效（用于登出）
func InvalidateTokenVersion(ctx context.Context, userID uuid.UUID) error {
	key := fmt.Sprintf("%s%s", tokenVersionKeyPrefix, userID.String())
	return config.RDB.Del(ctx, key).Err()
}
