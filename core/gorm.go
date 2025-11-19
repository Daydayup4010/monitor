package core

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"
	"uu/config"
	"uu/models"
)

func InitGorm() {
	if config.CONFIG.Mysql.Host == "" {
		config.Log.Warning("host is empty")
	}
	dsn := config.CONFIG.Mysql.Dsn()
	var mysqlLogger logger.Interface
	if config.CONFIG.Server.Env == "dev" {
		mysqlLogger = logger.Default.LogMode(logger.Info)
	} else {
		mysqlLogger = logger.Default.LogMode(logger.Error)
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: mysqlLogger, NamingStrategy: schema.NamingStrategy{SingularTable: true}})
	if err != nil {
		config.Log.Panicf("DB connect fail: %s", err)
	}
	err = db.AutoMigrate(&models.U{}, &models.BaseGoods{}, &models.User{}, &models.Settings{}, &models.APIKey{}, &models.Buff{}, &models.C5{}, &models.Steam{}) // migrate schema
	if err != nil {
		config.Log.Panicf("migrate schema fail: %s", err)
	}
	sqlDb, _ := db.DB()
	// SetMaxIdleConns: 设置空闲连接池中链接的最大数量
	sqlDb.SetMaxIdleConns(config.CONFIG.Mysql.MaxIdleConns)
	// SetMaxOpenConns: 设置打开数据库链接的最大数量
	sqlDb.SetMaxOpenConns(config.CONFIG.Mysql.MaxOpenConns)
	// SetConnMaxLifetime: 设置链接可复用的最大时间 (不要大于gin框架的timeout)
	sqlDb.SetConnMaxLifetime(10 * time.Second)
	config.DB = db
	config.Log.Info("gorm init success!")
}
