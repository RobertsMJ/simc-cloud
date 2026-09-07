package main

import (
	"github.com/RobertsMJ/simc-cloud/backend/config"
)

type SimConfig struct {
	*config.Config
	resultsQueueURL string
}

func NewSimConfig(cfg *config.Config) SimConfig {
	return SimConfig{
		Config: cfg,
		// resultsQueueURL: config.MustEnv("RESULTS_QUEUE_URL"),
	}
}
