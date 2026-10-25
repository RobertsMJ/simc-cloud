package admin

import "github.com/RobertsMJ/simc-cloud/backend/config"

type Config struct {
	Host string
	Port string
}

func NewConfig() *Config {
	return &Config{
		Host: config.OptionalEnv("ADMIN_HOST", "0.0.0.0"),
		Port: config.OptionalEnv("ADMIN_PORT", "8081"),
	}
}
