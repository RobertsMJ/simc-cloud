package main

import (
	"context"

	appconfig "github.com/RobertsMJ/simc-cloud-backend/config"
)

type Config struct {
	resultsQueueURL string
}

func LoadConfig(ctx context.Context) Config {
	return Config{
		resultsQueueURL: appconfig.MustEnv("RESULTS_QUEUE_URL"),
	}
}
