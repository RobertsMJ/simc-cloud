package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/RobertsMJ/simc-cloud-backend/internal/applog"
	"github.com/RobertsMJ/simc-cloud-backend/models"
	"github.com/RobertsMJ/simc-cloud-backend/sim"
	transport "github.com/RobertsMJ/simc-cloud-backend/transport/sqs"
	"github.com/aws/aws-lambda-go/lambda"
	sqssdk "github.com/aws/aws-sdk-go-v2/service/sqs"
)

var simulator sim.Simulator
var publisher *transport.Publisher[models.SimResult]

func init() {
	applog.Init()
	cfg := LoadConfig(context.Background())
	client := sqssdk.NewFromConfig(cfg.AWS)
	publisher = transport.NewPublisher[models.SimResult](context.Background(), client, cfg.resultsQueueName)
	simulator = sim.NewSimulator()
}

func handler(ctx context.Context, input models.SimRequest) (models.SimResult, error) {
	res, err := simulator.Run(ctx, &input)
	if err != nil {
		slog.Error("Simulation failed", "error", err)
		return models.SimResult{}, fmt.Errorf("simulation failed: %w", err)
	}

	publisher.Publish(ctx, res)

	return res, nil
}

func main() {
	lambda.Start(transport.NewRequestHandler(handler))
}
