package core

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"uu/config"
)

// InitConf 读取config yaml文件
func InitConf() {
	const ConfigFile = "settings.yaml"
	c := &config.Config{}
	yamlConf, err := os.ReadFile(ConfigFile)
	if err != nil {
		panic(fmt.Errorf("get Yaml Config file erro: %s", err))
		//global.LOG.Panicf("get Yaml Config file erro: %s", err)
	}
	err = yaml.Unmarshal(yamlConf, c)
	if err != nil {
		log.Fatalf("结构体映射错误: %s", err)
		//global.LOG.Panicf("结构体映射错误: %s", err)
	}
	//global.LOG.Info("config yaml load success!")
	log.Println("config yaml load success!")
	config.CONFIG = c
}
