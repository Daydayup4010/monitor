package main

import (
	"uu/config"
	"uu/core"
	"uu/models"
)

func main() {
	core.InitConf()
	core.InitLogger()
	core.InitGorm()
	core.InitRedis()
	models.InitKeys()
	//go services.UpdateBaseGoodsScheduler()
	//services.UpdateAllPlatformData()
	//go services.UpdateAllPlatformData()
	r := core.InitRouter()
	addr := config.CONFIG.Server.GetAddr()
	err := r.Run(addr)
	if err != nil {
		config.Log.Panicf("server start fail: %s", err)
	}

}
