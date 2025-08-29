package config

import "fmt"

type Server struct {
	Host   string `yaml:"host"`
	Port   string `yaml:"port"`
	Env    string `yaml:"env"`
	JwtKey string `yaml:"jwtKey"`
}

func (s *Server) GetAddr() string {
	return fmt.Sprintf("%s:%s", s.Host, s.Port)
}
