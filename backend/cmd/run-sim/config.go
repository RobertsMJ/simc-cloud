package main

import (
	"context"

	appconfig "github.com/RobertsMJ/simc-cloud-backend/config"
)

type Config struct {
	resultsQueueName string
}

func LoadConfig(ctx context.Context) Config {
	return Config{
		resultsQueueName: appconfig.MustEnv("RESULTS_QUEUE_NAME"),
	}
}
