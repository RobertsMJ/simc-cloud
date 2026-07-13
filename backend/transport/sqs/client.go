package sqs

import (
	"context"
	"log/slog"

	"github.com/RobertsMJ/simc-cloud-backend/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func NewClient(ctx context.Context, cfg aws.Config) *sqs.Client {
	if local := config.OptionalEnv("USE_LOCAL", "false"); local == "true" {
		slog.Debug("Using local SQS endpoint for testing")
		// cfg.Region = "elasticmq"
		cfg.BaseEndpoint = aws.String("http://elasticmq:9324")
		// cfg.Credentials = aws.NewCredentialsCache(aws.AnonymousCredentials{})
		return sqs.NewFromConfig(cfg)
	}
	return sqs.NewFromConfig(config.LoadAWS(ctx))
}
