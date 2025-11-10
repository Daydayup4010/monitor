package models

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"
	"uu/config"
	"uu/utils"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/scrypt"
	"gorm.io/gorm"
)

const (
	RoleNormal int64 = 0
	RoleVip    int64 = 1
	RoleAdmin  int64 = 2
)

type User struct {
	ID           uuid.UUID      `json:"id" gorm:"type:char(36);primaryKey"`
	UserName     string         `json:"username" gorm:"type:varchar(255);Index"`
	Email        string         `json:"email" gorm:"type:varchar(255);uniqueIndex:idx_user_email_deleted"`
	Password     string         `gorm:"type:varchar(255)"`
	WechatOpenID *string        `json:"wechat_openid" gorm:"type:varchar(128);uniqueIndex:idx_wechat_openid_deleted;default:NULL"`
	Role         int64          `json:"role" gorm:"type:int;default:0"`
	VipExpiry    time.Time      `json:"vip_expiry" gorm:"type:datetime;default:CURRENT_TIMESTAMP"`
	LastLogin    time.Time      `json:"last_login" gorm:"type:datetime;default:CURRENT_TIMESTAMP"`
	CreatedAt    time.Time      `gorm:"type:datetime;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"type:datetime;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"type:datetime;index;uniqueIndex:idx_user_email_deleted;uniqueIndex:idx_wechat_openid_deleted" json:"deleted_at,omitempty"`
}

type UserResponse struct {
	Id        string    `json:"id"`
	UserName  string    `json:"user_name"`
	Email     string    `json:"email"`
	Role      int64     `json:"role"`
	VipExpiry time.Time `json:"vip_expiry,omitempty"`
	LastLogin time.Time `json:"last_login,omitempty"`
}

func CreateUser(user *User) int {
	if IfExistEmail(user.Email) {
		return utils.ErrCodeEmailTaken
	}
	err := config.DB.Create(user).Error
	if err != nil {
		config.Log.Errorf("Create User error: %v", err)
		return utils.ErrCodeCreateUser
	}
	return utils.SUCCESS
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

func UpdateUserName(name, id string) int {
	err := config.DB.Model(&User{}).Where("id = ?", id).Update("user_name", name).Error
	if err != nil {
		config.Log.Errorf("update user name error :%v", err)
		return utils.ErrCodeUpdateUser
	}
	return utils.SUCCESS
}

func ResetPassword(email, password string) int {
	if !IfExistEmail(email) {
		return utils.ErrCodeUserNotFound
	}
	err := config.DB.Model(&User{}).Where("email = ?", email).Update("password", password).Error
	if err != nil {
		config.Log.Errorf("reset password fail: %v", err)
		return utils.ErrCodeUpdateUser
	}
	return utils.SUCCESS
}

func DeleteUser(id string) int {
	result := config.DB.Where("id = ?", id).Delete(&User{})
	if result.RowsAffected == 0 {
		return utils.ErrCodeUserNotFound
	}
	if result.Error != nil {
		config.Log.Errorf("delete user fail: %v", result.Error)
		return utils.ErrCodeDeleteUser
	}
	return utils.SUCCESS
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

func QueryUserByOpenID(openID string) *User {
	if openID == "" {
		return nil
	}
	var user User
	err := config.DB.Where("wechat_open_id = ?", openID).First(&user).Error
	if err != nil {
		return nil
	}
	return &user
}

func CreateWechatUser(user *User) int {
	err := config.DB.Create(user).Error
	if err != nil {
		config.Log.Errorf("Create Wechat User error: %v", err)
		return utils.ErrCodeCreateUser
	}
	return utils.SUCCESS
}

func GetUserById(id string) (*User, int) {
	var user User
	err := config.DB.Where("id = ?", id).First(&user).Error
	if err != nil {
		config.Log.Errorf("query user by id error: %v", err)
		return nil, utils.ErrCodeUserNotFound
	}
	return &user, utils.SUCCESS
}

func UpdateUserLastLogin(user *User) {
	err := config.DB.Model(&User{}).Where("email = ?", user.Email).Update("last_login", user.LastLogin).Error
	if err != nil {
		config.Log.Errorf("update user last login fail: %v", err)
	}
}

func GetUserList(pageSize, pageNum int, search string) ([]UserResponse, int64, int) {
	var users []UserResponse
	var total int64

	query := config.DB.Model(&User{})

	// 添加搜索条件
	if search != "" {
		query = query.Where("user_name LIKE ? OR email LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	query.Count(&total)
	err := query.Select("id, user_name, email, role, vip_expiry, last_login").
		Limit(pageSize).
		Offset((pageNum - 1) * pageSize).
		Find(&users).Error

	if err != nil {
		config.Log.Errorf("Get user list error: %v", err)
		return users, total, utils.ErrCodeGetUserList
	}
	return users, total, utils.SUCCESS
}

func SaveEmailCode(email, code string, c context.Context) int {
	key := fmt.Sprintf("verify:%s", email)
	err := config.RDB.SetEx(c, key, code, 10*time.Minute).Err()
	if err != nil {
		config.Log.Errorf("generate verify code fail: %v", err)
		return utils.ErrCodeEmailCodeGenerate
	}
	return utils.SUCCESS
}

func VerifyEmailCode(email, code string, c context.Context) (bool, int) {
	key := fmt.Sprintf("verify:%s", email)
	saveCode, err := config.RDB.Get(c, key).Result()
	if err == redis.Nil {
		return false, utils.ErrCodeInvalidEmailCode
	} else if err != nil {
		return false, utils.ErrCodeInvalidEmailCode
	}
	if saveCode != code {
		return false, utils.ErrCodeInvalidEmailCode
	}
	config.RDB.Del(c, key)
	return true, utils.SUCCESS
}

func RenewVIP(userID string, days int) (time.Time, int) {
	// 获取当前VIP过期时间
	var user User
	var newExpiry time.Time
	if err := config.DB.First(&user, "id = ?", userID).Error; err != nil {
		config.Log.Errorf("query user fail: %v", err)
		return newExpiry, utils.ErrCodeUserNotFound
	}

	// 计算新的过期时间（days参数表示月数）
	months := days
	if user.VipExpiry.IsZero() || user.VipExpiry.Before(time.Now()) {
		// 如果VIP已过期或未设置，从当前时间开始算
		newExpiry = time.Now().AddDate(0, months, 0)
	} else {
		// 如果VIP未过期，在原有基础上续费
		newExpiry = user.VipExpiry.AddDate(0, months, 0)
	}

	// 更新VIP到期时间和角色
	updates := map[string]interface{}{
		"vip_expiry": newExpiry,
		"role":       RoleVip, // 设置为VIP角色
	}

	err := config.DB.Model(&user).Updates(updates).Error
	if err != nil {
		config.Log.Errorf("update vip expiry fail: %v", err)
		return newExpiry, utils.ErrCodeUpdateUser
	}
	return newExpiry, utils.SUCCESS
}
