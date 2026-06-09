package main

import (
	"context"
	"log/slog"

	"github.com/RobertsMJ/simc-cloud-backend/db"
	"github.com/RobertsMJ/simc-cloud-backend/internal/applog"
	"github.com/RobertsMJ/simc-cloud-backend/job"
	"github.com/RobertsMJ/simc-cloud-backend/models"
	transport "github.com/RobertsMJ/simc-cloud-backend/transport/sqs"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var writer job.ResultWriter

func init() {
	applog.Init()
	cfg := LoadConfig(context.Background())
	client := dynamodb.NewFromConfig(cfg.AWS)
	writer = db.NewJobRepository(db.JobRepositoryConfig{
		Client:    client,
		TableName: cfg.tableName,
	})
}

func handler(ctx context.Context, input models.SimResult) error {
	if err := writer.WriteResult(ctx, input); err != nil {
		slog.Error("Failed to write result", slog.Any("error", err))
		return err
	}
	return nil
}

func main() {
	lambda.Start(transport.NewConsumer(handler))
}
