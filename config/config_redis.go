package config

import "fmt"

type Redis struct {
	Host     string `yaml:"host"`
	Port     string `json:"port"`
	Password string `json:"password"`
	Db       int    `json:"db"`
}

func (r *Redis) Addr() string {
	return fmt.Sprintf("%s:%s", r.Host, r.Port)
}
