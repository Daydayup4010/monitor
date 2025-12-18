package config

import (
	"uu/utils"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Config struct {
	Mysql   *Mysql              `yaml:"mysql"`
	Logger  *Logger             `yaml:"logger"`
	Server  *Server             `yaml:"server"`
	Redis   *Redis              `yaml:"redis"`
	Email   *utils.EmailService `yaml:"email"`
	Wechat  *Wechat             `yaml:"wechat"`
	SteamDt *SteamDt            `yaml:"steamDt"`
	Payment *Payment            `yaml:"payment"`
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

var (
	CONFIG *Config
	Log    *logrus.Logger
	DB     *gorm.DB
	RDB    *redis.Client
)
