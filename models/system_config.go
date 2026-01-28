package models

import (
	"uu/config"
)

// SystemConfig 系统配置表
type SystemConfig struct {
	ID    uint   `gorm:"primaryKey;autoIncrement"`
	Key   string `gorm:"type:varchar(100);uniqueIndex;not null"`
	Value string `gorm:"type:varchar(500)"`
	Desc  string `gorm:"type:varchar(255)"` // 配置描述
}

const (
	// ConfigKeyMinAppVipEnabled 小程序VIP开通入口开关
	ConfigKeyMinAppVipEnabled = "minapp_vip_enabled"
)

// GetSystemConfig 获取系统配置
func GetSystemConfig(key string) (string, bool) {
	var cfg SystemConfig
	err := config.DB.Where("`key` = ?", key).First(&cfg).Error
	if err != nil {
		return "", false
	}
	return cfg.Value, true
}

// SetSystemConfig 设置系统配置
func SetSystemConfig(key, value, desc string) error {
	var cfg SystemConfig
	err := config.DB.Where("`key` = ?", key).First(&cfg).Error
	if err != nil {
		// 不存在则创建
		cfg = SystemConfig{
			Key:   key,
			Value: value,
			Desc:  desc,
		}
		return config.DB.Create(&cfg).Error
	}
	// 存在则更新
	cfg.Value = value
	if desc != "" {
		cfg.Desc = desc
	}
	return config.DB.Save(&cfg).Error
}

// GetMinAppConfig 获取小程序配置（公开API）
func GetMinAppConfig() map[string]interface{} {
	result := make(map[string]interface{})

	// VIP开通入口开关，默认关闭
	vipEnabled, exists := GetSystemConfig(ConfigKeyMinAppVipEnabled)
	if !exists || vipEnabled != "1" {
		result["vip_enabled"] = false
	} else {
		result["vip_enabled"] = true
	}

	return result
}

// InitSystemConfig 初始化系统配置默认值
func InitSystemConfig() {
	// 小程序VIP开关，默认关闭
	if _, exists := GetSystemConfig(ConfigKeyMinAppVipEnabled); !exists {
		SetSystemConfig(ConfigKeyMinAppVipEnabled, "0", "小程序VIP开通入口开关(0=关闭,1=开启)")
	}
}
