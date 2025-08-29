package config

type Mysql struct {
	Host         string `yaml:"host"`
	Port         string `yaml:"port"`
	Config       string `yaml:"config"`
	DB           string `yaml:"db"`
	Username     string `yaml:"username"`
	Password     string `yaml:"password"`
	LogLevel     string `yaml:"log_level"`
	MaxIdleConns int    `yaml:"maxIdleConns"`
	MaxOpenConns int    `yaml:"maxOpenConns"`
}

func (m *Mysql) Dsn() string {
	return m.Username + ":" + m.Password + "@tcp(" + m.Host + ":" + m.Port + ")/" + m.DB + "?" + m.Config
}
