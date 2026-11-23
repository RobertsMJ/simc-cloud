package scyllax

import "github.com/RobertsMJ/simc-cloud/backend/config"

type ScyllaConfig struct {
	Hosts            []string `env:"SCYLLA_HOSTS,required"`
	Keyspace         string   `env:"SCYLLA_KEYSPACE,required"`
	Username         string   `env:"SCYLLA_USERNAME,required"`
	Password         string   `env:"SCYLLA_PASSWORD,required"`
	DisableDiscovery bool     `env:"SCYLLA_DISABLE_DISCOVERY"`
}

func NewScyllaConfig() (*ScyllaConfig, error) {
	return config.Parse[ScyllaConfig]("")
}
