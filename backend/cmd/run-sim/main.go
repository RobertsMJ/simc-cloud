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
)

type resultPublisher interface {
	Publish(ctx context.Context, result models.SimResult) error
}

var simulator sim.Simulator
var publisher resultPublisher

func init() {
	applog.Init()
	cfg := LoadConfig(context.Background())
	publisher = transport.NewPublisher[models.SimResult](context.Background(), cfg.resultsQueueName)
	simulator = sim.NewSimulator()
}

func handler(ctx context.Context, input models.SimRequest) error {
	res, err := simulator.Run(ctx, &input)
	if err != nil {
		slog.Error("Simulation failed", "error", err)
		return fmt.Errorf("simulation failed: %w", err)
	}

	if err := publisher.Publish(ctx, res); err != nil {
		slog.Error("Failed to publish result", "error", err)
		return fmt.Errorf("failed to publish result: %w", err)
	}

	return nil
}

func main() {
	lambda.Start(transport.NewConsumer(handler))
}
