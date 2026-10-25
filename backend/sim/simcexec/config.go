package simcexec

import (
	"runtime"

	"github.com/RobertsMJ/simc-cloud/backend/config"
)

type Config struct {
	BinPath string
	Threads int
}

func NewConfig() *Config {
	bin := config.OptionalEnv("SIM_BINARY_PATH", "/app/simc")
	return &Config{
		BinPath: bin,
		Threads: runtime.GOMAXPROCS(0),
	}
}
