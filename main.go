package main

import (
	"uu/config"
	"uu/core"
	"uu/models"
	"uu/services"
)

func main() {
	core.InitConf()
	core.InitLogger()
	core.InitGorm()
	core.InitRedis()
	models.InitKeys()
	go services.UpdateBaseGoodsScheduler()
	go services.UpdateAllGoodsScheduler()
	go services.UpdateIconScheduler()
	go services.StartDailyPriceHistoryScheduler()
	r := core.InitRouter()
	addr := config.CONFIG.Server.GetAddr()
	err := r.Run(addr)
	if err != nil {
		config.Log.Panicf("server start fail: %s", err)
	}

}
