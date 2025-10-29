package utils

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
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
	Accept            string `yaml:"accept"`
	AcceptLanguage    string `yaml:"accept-language"`
	AcceptEncoding    string `yaml:"accept-encoding"`
	UserAgent         string `yaml:"user-agent"`
	TimezoneOffset    string `yaml:"timezone-offset"`
	DeviceName        string `yaml:"device-name"`
	DeviceIdWeak      string `yaml:"device-id-weak"`
	ScreenScale       string `yaml:"screen-scale"`
	Resolution        string `yaml:"resolution"`
	Locale            string `yaml:"locale"`
	DeviceId          string `yaml:"device-id"`
	Connection        string `yaml:"connection"`
	LocaleSupported   string `yaml:"locale-supported"`
	Timezone          string `yaml:"timezone"`
	Network           string `yaml:"network"`
	Product           string `yaml:"product"`
	TimezoneOffsetDst string `yaml:"timezone-offset-dst"`
	Model             string `yaml:"model"`
	AppVersion        string `yaml:"app-version"`
	ScreenSize        string `yaml:"screen-size"`
	AppVersionCode    string `yaml:"app-version-code"`
	SystemVersion     string `yaml:"system-version"`
	Session           string
	Csrf              string
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
