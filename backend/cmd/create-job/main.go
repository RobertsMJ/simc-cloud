package main

import (
	"context"

	"github.com/RobertsMJ/simc-cloud-backend/logger"
	"github.com/RobertsMJ/simc-cloud-backend/models"
	transport "github.com/RobertsMJ/simc-cloud-backend/transport/apigw"
	"github.com/aws/aws-lambda-go/lambda"
)

func init() {
	logger.LoadLogger()
}

func handler(ctx context.Context, req models.CreateJobRequest) (models.CreateJobResponse, error) {
	logger.Error("CreateJob handler not implemented")
	return models.CreateJobResponse{}, models.ErrNotImplemented
}

func main() {
	lambda.Start(transport.NewRequestHandler(handler))
}
