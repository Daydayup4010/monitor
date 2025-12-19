package config

import (
	"uu/utils"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Config struct {
	Mysql      *Mysql              `yaml:"mysql"`
	Logger     *Logger             `yaml:"logger"`
	Server     *Server             `yaml:"server"`
	Redis      *Redis              `yaml:"redis"`
	Email      *utils.EmailService `yaml:"email"`
	Wechat     *Wechat             `yaml:"wechat"`
	SteamDt    *SteamDt            `yaml:"steamDt"`
	Payment    *Payment            `yaml:"payment"`
	ErrorAlert *ErrorAlert         `yaml:"errorAlert"`
}

type Payment struct {
	MchId     string `yaml:"mch_id"`     // YunGouOS商户号
	ApiKey    string `yaml:"api_key"`    // API密钥
	NotifyUrl string `yaml:"notify_url"` // 支付回调地址
}

type Wechat struct {
	AppID     string `yaml:"app_id"`
	AppSecret string `yaml:"app_secret"`
}

type SteamDt struct {
	Key string `yaml:"key"`
}

// ErrorAlert 错误告警配置
type ErrorAlert struct {
	Enabled     bool     `yaml:"enabled"`      // 是否启用错误告警
	Recipients  []string `yaml:"recipients"`   // 接收告警邮件的邮箱列表
	MinLevel    string   `yaml:"min_level"`    // 最低告警级别: error, fatal, panic
	RateLimit   int      `yaml:"rate_limit"`   // 每分钟最多发送的告警邮件数量
	Cooldown    int      `yaml:"cooldown"`     // 相同错误的冷却时间（秒）
	BatchWindow int      `yaml:"batch_window"` // 批量发送窗口（秒），0表示立即发送
}

var (
	CONFIG *Config
	Log    *logrus.Logger
	DB     *gorm.DB
	RDB    *redis.Client
)
