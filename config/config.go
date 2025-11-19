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
