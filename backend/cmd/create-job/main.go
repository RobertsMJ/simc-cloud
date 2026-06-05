package main

import (
	"context"
	"log/slog"

	"github.com/RobertsMJ/simc-cloud-backend/internal/applog"
	"github.com/RobertsMJ/simc-cloud-backend/models"
	transport "github.com/RobertsMJ/simc-cloud-backend/transport/apigw"
	"github.com/aws/aws-lambda-go/lambda"
)

func init() {
	applog.Init()
}

func handler(ctx context.Context, req models.CreateJobRequest) (models.CreateJobResponse, error) {
	slog.Error("CreateJob handler not implemented")
	return models.CreateJobResponse{}, models.ErrNotImplemented
}

func main() {
	lambda.Start(transport.NewRequestHandler(handler))
}
