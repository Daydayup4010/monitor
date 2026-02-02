package models

import (
	"context"
	"fmt"
	"time"
	"uu/config"

	"github.com/google/uuid"
)

const tokenVersionKeyPrefix = "user_token_version:"

// 客户端类型常量
const (
	ClientTypeWeb         = "web"
	ClientTypeMiniprogram = "miniprogram"
)

// GenerateTokenVersion 生成新的 token 版本号（使用时间戳）
func GenerateTokenVersion() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// SetTokenVersion 存储用户的当前 token 版本号（按客户端类型区分）
func SetTokenVersion(ctx context.Context, userID uuid.UUID, clientType string, version string) error {
	key := fmt.Sprintf("%s%s:%s", tokenVersionKeyPrefix, clientType, userID.String())
	// 设置 25 小时过期（比 token 有效期稍长）
	return config.RDB.Set(ctx, key, version, 25*time.Hour).Err()
}

// GetTokenVersion 获取用户当前的 token 版本号（按客户端类型区分）
func GetTokenVersion(ctx context.Context, userID uuid.UUID, clientType string) (string, error) {
	key := fmt.Sprintf("%s%s:%s", tokenVersionKeyPrefix, clientType, userID.String())
	return config.RDB.Get(ctx, key).Result()
}

// ValidateTokenVersion 验证 token 版本号是否有效（按客户端类型区分）
func ValidateTokenVersion(ctx context.Context, userID uuid.UUID, clientType string, version string) bool {
	currentVersion, err := GetTokenVersion(ctx, userID, clientType)
	if err != nil {
		return false
	}
	return currentVersion == version
}

// InvalidateTokenVersion 使 token 版本失效（用于登出，按客户端类型区分）
func InvalidateTokenVersion(ctx context.Context, userID uuid.UUID, clientType string) error {
	key := fmt.Sprintf("%s%s:%s", tokenVersionKeyPrefix, clientType, userID.String())
	return config.RDB.Del(ctx, key).Err()
}

// InvalidateAllTokenVersions 使所有客户端的 token 版本失效
func InvalidateAllTokenVersions(ctx context.Context, userID uuid.UUID) error {
	// 删除 web 和 miniprogram 的 token
	webKey := fmt.Sprintf("%s%s:%s", tokenVersionKeyPrefix, ClientTypeWeb, userID.String())
	miniprogramKey := fmt.Sprintf("%s%s:%s", tokenVersionKeyPrefix, ClientTypeMiniprogram, userID.String())
	return config.RDB.Del(ctx, webKey, miniprogramKey).Err()
}
