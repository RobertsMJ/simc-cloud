package main

import (
	"context"

	appconfig "github.com/RobertsMJ/simc-cloud-backend/config"
	"github.com/aws/aws-sdk-go-v2/aws"
)

type Config struct {
	AWS              aws.Config
	resultsQueueName string
}

func LoadConfig(ctx context.Context) Config {
	return Config{
		AWS:              appconfig.LoadAWS(ctx),
		resultsQueueName: appconfig.MustEnv("RESULTS_QUEUE_NAME"),
	}
}
