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
	ID        uuid.UUID      `gorm:"type:char(36);primaryKey"`
	UserName  string         `json:"username" gorm:"type:varchar(255);Index"`
	Email     string         `json:"email" gorm:"type:varchar(255);uniqueIndex"`
	Password  string         `gorm:"type:varchar(255)"`
	Role      int64          `json:"role" gorm:"type:int;default:0"`
	VipExpiry time.Time      `json:"vip_expiry" gorm:"type:datetime;default:CURRENT_TIMESTAMP"`
	LastLogin time.Time      `json:"last_login" gorm:"type:datetime;default:CURRENT_TIMESTAMP"`
	CreatedAt time.Time      `gorm:"type:datetime;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time      `gorm:"type:datetime;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"type:datetime;index" json:"deleted_at,omitempty"`
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

func IfExistUser(id uuid.UUID) bool {
	var count int64
	config.DB.Model(&User{}).Where("id = ?", id).Count(&count)
	if count > 0 {
		return true
	}
	return false
}

func UpdateUserName(name, id string) error {
	err := config.DB.Model(&User{}).Where("id = ?", id).Update("username", name).Error
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

// IsValidVIP 检查VIP是否有效
func (u *User) IsValidVIP() bool {
	return u.Role == RoleVip && (u.VipExpiry.After(time.Now()))
}

// CanAccessVIPContent 检查VIP访问权限
func (u *User) CanAccessVIPContent() bool {
	return u.IsValidVIP() || u.Role == RoleAdmin
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
	err := config.DB.Model(&User{}).Select("user_name,email,role,vip_expiry,last_login").Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&users).Error
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
