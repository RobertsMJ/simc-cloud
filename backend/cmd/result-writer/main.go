package main

import (
	"context"
	"log/slog"

	"github.com/RobertsMJ/simc-cloud-backend/internal/applog"
	"github.com/RobertsMJ/simc-cloud-backend/job"
	"github.com/RobertsMJ/simc-cloud-backend/models"
	transport "github.com/RobertsMJ/simc-cloud-backend/transport/sqs"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var client *dynamodb.Client
var writer *job.ResultWriter

func init() {
	applog.Init()
	cfg := LoadConfig(context.Background())
	client = dynamodb.NewFromConfig(cfg.AWS)
}

func handler(ctx context.Context, input models.SimResult) (models.SimResult, error) {
	if err := (*writer).WriteResult(ctx, input); err != nil {
		slog.Error("Failed to write result", slog.Any("error", err))
		return models.SimResult{}, err
	}
	return input, nil
}

func main() {
	lambda.Start(transport.NewRequestHandler(handler))
}
