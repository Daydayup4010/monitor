package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

var JWTSecret = []byte("asdfghjkzxcvbnm") // 生产环境应从环境变量获取

// Claims 自定义JWT声明
type Claims struct {
	UserID       uuid.UUID `json:"user_id"`
	Username     string    `json:"username"`
	Role         int64     `json:"role"`
	VipExpiry    time.Time `json:"vip_expiry"`
	Email        string    `json:"email"`
	TokenVersion string    `json:"token_version"` // 用于单设备登录验证
	jwt.RegisteredClaims
}

// GenerateJWT 生成JWT令牌
func GenerateJWT(userID uuid.UUID, username string, role int64, vipExpiry time.Time, email string, tokenVersion string) (string, error) {
	claims := &Claims{
		UserID:       userID,
		Username:     username,
		Role:         role,
		VipExpiry:    vipExpiry,
		Email:        email,
		TokenVersion: tokenVersion,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 24小时有效期
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "monitor",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JWTSecret)
}

// ParseJWT 解析和验证JWT令牌
func ParseJWT(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return JWTSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}
