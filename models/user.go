package models

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/scrypt"
	"gorm.io/gorm"
	"time"
	"uu/config"
)

const (
	RoleNormal int64 = 0
	RoleVip    int64 = 1
	RoleAdmin  int64 = 2
)

type User struct {
	ID        uuid.UUID      `json:"id" gorm:"type:char(36);primaryKey"`
	UserName  string         `json:"username" gorm:"type:varchar(255);Index"`
	Email     string         `json:"email" gorm:"type:varchar(255);uniqueIndex:idx_user_email_deleted"`
	Password  string         `gorm:"type:varchar(255)"`
	Role      int64          `json:"role" gorm:"type:int;default:0"`
	VipExpiry time.Time      `json:"vip_expiry" gorm:"type:datetime;default:CURRENT_TIMESTAMP"`
	LastLogin time.Time      `json:"last_login" gorm:"type:datetime;default:CURRENT_TIMESTAMP"`
	CreatedAt time.Time      `gorm:"type:datetime;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time      `gorm:"type:datetime;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"type:datetime;index;uniqueIndex:idx_user_email_deleted" json:"deleted_at,omitempty"`
}

type UserResponse struct {
	Id        string    `json:"id"`
	UserName  string    `json:"user_name"`
	Email     string    `json:"email"`
	Role      int64     `json:"role"`
	VipExpiry time.Time `json:"vip_expiry,omitempty"`
	LastLogin time.Time `json:"last_login,omitempty"`
}

func CreateUser(user *User) error {
	if IfExistEmail(user.Email) {
		return fmt.Errorf("email registered")
	}
	err := config.DB.Create(user).Error
	return err
}

func IfExistEmail(email string) bool {
	var count int64
	config.DB.Model(&User{}).Where("email = ?", email).Count(&count)
	if count > 0 {
		return true
	}
	return false
}

func IfExistUser(name string) bool {
	var count int64
	config.DB.Model(&User{}).Where("user_name = ?", name).Count(&count)
	if count > 0 {
		return true
	}
	return false
}

func UpdateUserName(name, id string) error {
	err := config.DB.Model(&User{}).Where("id = ?", id).Update("user_name", name).Error
	return err
}

func ResetPassword(email, password string) error {
	if !IfExistEmail(email) {
		return fmt.Errorf("email not exist")
	}
	err := config.DB.Model(&User{}).Where("email = ?", email).Update("password", password).Error
	return err
}

func DeleteUser(id string) error {
	result := config.DB.Where("id = ?", id).Delete(&User{})
	if result.RowsAffected == 0 {
		return fmt.Errorf("user not exist")
	}
	return result.Error
}

func IsValidVIP(role int64, expiry time.Time) bool {
	return role == RoleVip && (expiry.After(time.Now()))
}

func CanAccessVIPContent(role int64, expiry time.Time) bool {
	return IsValidVIP(role, expiry) || role == RoleAdmin
}

// ScryptPw 密码加密
func ScryptPw(password string) string {
	const KeyLen = 10
	salt := make([]byte, 8)
	salt = []byte{12, 33, 11, 51, 111, 200, 255, 10}
	key, err := scrypt.Key([]byte(password), salt, 16384, 8, 1, KeyLen)
	if err != nil {
		config.Log.Warningf("password scrypt fail: %v", err)
	}
	pw := base64.StdEncoding.EncodeToString(key) // key是byte数组，最后转换为字符串
	return pw
}

func QueryUser(email string) *User {
	var user User
	err := config.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		config.Log.Errorf("query user error: %v", err)
		return nil
	}
	return &user
}

func UpdateUserLastLogin(user *User) {
	err := config.DB.Model(&User{}).Where("email = ?", user.Email).Update("last_login", user.LastLogin).Error
	if err != nil {
		config.Log.Errorf("update user last login fail: %v", err)
	}
}

func GetUserList(pageSize, pageNum int) ([]UserResponse, int64) {
	var users []UserResponse
	var total int64
	config.DB.Model(&User{}).Count(&total)
	err := config.DB.Model(&User{}).Select("id, user_name,email,role,vip_expiry,last_login").Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&users).Error
	if err != nil {
		config.Log.Errorf("Get user list error: %v", err)
	}
	return users, total
}

func SaveEmailCode(email, code string, c context.Context) error {
	key := fmt.Sprintf("verify:%s", email)
	err := config.RDB.SetEx(c, key, code, 10*time.Minute).Err()
	return err
}

func VerifyEmailCode(email, code string, c context.Context) (bool, error) {
	key := fmt.Sprintf("verify:%s", email)
	saveCode, err := config.RDB.Get(c, key).Result()
	if err == redis.Nil {
		return false, fmt.Errorf("code not exist")
	} else if err != nil {
		return false, err
	}
	if saveCode != code {
		return false, nil
	}
	config.RDB.Del(c, key)
	return true, nil
}

func RenewVIP(userID string, days int) (time.Time, error) {
	// 获取当前VIP过期时间
	var user User
	var newExpiry time.Time
	if err := config.DB.First(&user, "id = ?", userID).Error; err != nil {
		config.Log.Errorf("query user fail: %v", err)
		return newExpiry, fmt.Errorf("query user fail")
	}

	// 计算新的过期时间
	if user.VipExpiry.IsZero() || user.VipExpiry.Before(time.Now()) {
		// 如果VIP已过期或未设置，从当前时间开始算
		newExpiry = time.Now().AddDate(0, 0, days)
	} else {
		// 如果VIP未过期，在原有基础上续费
		newExpiry = user.VipExpiry.AddDate(0, 0, days)
	}

	err := config.DB.Model(&user).Update("vip_expiry", newExpiry).Error
	if err != nil {
		config.Log.Errorf("update vip expiry fail: %v", err)
		return newExpiry, fmt.Errorf("update vip expiry fail")
	}
	return newExpiry, nil
}
