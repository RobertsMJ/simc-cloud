package job

import "github.com/RobertsMJ/simc-cloud/backend/config"

type Config struct {
	MaxGearsets int `env:"MAX_GEARSETS" envDefault:"100"`
}

func NewConfig() (*Config, error) {
	return config.Parse[Config]("APP_")
}
