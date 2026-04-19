package main

import (
	"context"
	"fmt"

	"github.com/RobertsMJ/simc-cloud-backend/logger"
	"github.com/RobertsMJ/simc-cloud-backend/sim"
	"github.com/RobertsMJ/simc-cloud-backend/simc"
	transport "github.com/RobertsMJ/simc-cloud-backend/transport/sqs"
	"github.com/aws/aws-lambda-go/lambda"
)

func init() {
	logger.LoadLogger()
}

func handler(ctx context.Context, input *simc.Input) (*simc.Output, error) {

	logger.Info(fmt.Sprintf("Starting simulation:\n%v", input), input)
	simulator := sim.NewSimulator()
	result, err := simulator.Run(ctx, input)
	if err != nil {
		logger.Error("Simulation failed", "error", err)
		return nil, fmt.Errorf("simulation failed: %w", err)
	}

	return &result, nil
}

func main() {
	lambda.Start(transport.NewRequestHandler(handler))
}
