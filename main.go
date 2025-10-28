package main

import (
	"uu/config"
	"uu/core"
	"uu/services"
)

func main() {
	core.InitConf()
	core.InitLogger()
	core.InitGorm()
	core.InitRedis()
	go services.StartBuffFullUpdateScheduler()
	go services.StartUUFullUpdateScheduler()
	go services.StartVerifyToken()
	r := core.InitRouter()
	addr := config.CONFIG.Server.GetAddr()
	err := r.Run(addr)
	if err != nil {
		config.Log.Panicf("server start fail: %s", err)
	}

}
