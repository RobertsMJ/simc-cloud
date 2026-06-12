package main

import (
	"context"

	"github.com/RobertsMJ/simc-cloud-backend/config"
	"github.com/aws/aws-sdk-go-v2/aws"
)

type Config struct {
	AWS       aws.Config
	tableName string
}

func LoadConfig(ctx context.Context) Config {
	return Config{
		AWS:       config.LoadAWS(ctx),
		tableName: config.MustEnv("SIM_TABLE_NAME"),
	}
}
