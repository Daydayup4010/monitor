package main

import (
	"uu/config"
	"uu/core"
)

func main() {
	core.InitConf()
	core.InitLogger()
	core.InitGorm()
	core.InitRedis()
	r := core.InitRouter()
	addr := config.CONFIG.Server.GetAddr()
	err := r.Run(addr)
	if err != nil {
		config.Log.Panicf("server start fail: %s", err)
	}
}
