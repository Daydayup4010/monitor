package utils

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type Headers struct {
	UU   *UU   `yaml:"uu"`
	Buff *Buff `yaml:"buff"`
}

type UU struct {
	Accept         string `yaml:"accept"`
	AcceptLanguage string `yaml:"accept-language"`
	AppVersion     string `yaml:"app-version"`
	AppType        string `yaml:"apptype"`
	ContentType    string `yaml:"content-type"`
	Platform       string `yaml:"platform"`
	UserAgent      string `yaml:"user-agent"`
	Uk             string
	Authorization  string
}

type Buff struct {
	Accept         string `yaml:"accept"`
	AcceptLanguage string `yaml:"accept-language"`
	ContentType    string `yaml:"content-type"`
	UserAgent      string `yaml:"user-agent"`
	Session        string
	Csrf           string
}

func ReadeHeaders() *Headers {
	c := &Headers{}
	const ConfigFile = "headers.yaml"
	yamlConf, err := os.ReadFile(ConfigFile)
	if err != nil {
		panic(fmt.Errorf("get Headers yaml file erro: %s", err))
		//global.LOG.Panicf("get Yaml Config file erro: %s", err)
	}
	err = yaml.Unmarshal(yamlConf, c)
	if err != nil {
		log.Fatalf("结构体映射错误: %s", err)
		//global.LOG.Panicf("结构体映射错误: %s", err)
	}
	//global.LOG.Info("config yaml load success!")
	log.Println("header yaml load success!")
	return c
}

var Header = ReadeHeaders()
