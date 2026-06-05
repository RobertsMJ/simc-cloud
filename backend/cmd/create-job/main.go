package main

import (
	"context"
	"log/slog"
	"os"
	"strings"

	"github.com/RobertsMJ/simc-cloud-backend/models"
	transport "github.com/RobertsMJ/simc-cloud-backend/transport/apigw"
	"github.com/aws/aws-lambda-go/lambda"
)

func init() {
	initLogger()
}

func initLogger() {
	var level slog.Level
	switch strings.ToLower(os.Getenv("LOG_LEVEL")) {
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})))
}

func handler(ctx context.Context, req models.CreateJobRequest) (models.CreateJobResponse, error) {
	slog.Error("CreateJob handler not implemented")
	return models.CreateJobResponse{}, models.ErrNotImplemented
}

func main() {
	lambda.Start(transport.NewRequestHandler(handler))
}
