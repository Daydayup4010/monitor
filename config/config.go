package config

import (
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Config struct {
	Mysql  *Mysql  `yaml:"mysql"`
	Logger *Logger `yaml:"logger"`
	Server *Server `yaml:"server"`
	Redis  *Redis  `yaml:"redis"`
}

var (
	CONFIG *Config
	Log    *logrus.Logger
	DB     *gorm.DB
	RDB    *redis.Client
)
