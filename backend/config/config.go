package config

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
)

func LoadAWS(ctx context.Context) aws.Config {
	cfg, err := awsconfig.LoadDefaultConfig(ctx)
	if err != nil {
		panic(fmt.Errorf("failed to load AWS config: %w", err))
	}
	return cfg
}

func MustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		panic(fmt.Sprintf("required environment variable %q is not set", key))
	}
	return v
}

func OptionalEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
