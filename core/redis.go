package core

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
	"uu/config"
)

func InitRedis() {
	if config.CONFIG.Redis.Host == "" {
		config.Log.Warning("host is empty")
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.CONFIG.Redis.Addr(),
		Password: config.CONFIG.Redis.Password,
		DB:       config.CONFIG.Redis.Db,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if _, err := rdb.Ping(ctx).Result(); err != nil {
		config.Log.Panicf("redis connet fail: %s", err)
	}
	config.RDB = rdb
	config.Log.Info("redis init success")
}
