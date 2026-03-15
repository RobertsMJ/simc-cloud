package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/RobertsMJ/simc-cloud-backend/logger"
	"github.com/RobertsMJ/simc-cloud-backend/models"
	"github.com/RobertsMJ/simc-cloud-backend/sim"
	transport "github.com/RobertsMJ/simc-cloud-backend/transport/apigw"
	"github.com/aws/aws-lambda-go/lambda"
)

func init() {
	logger.LoadLogger()
}

func handler(ctx context.Context, event transport.Request) (transport.Response, error) {
	start := time.Now()
	defer func() {
		logger.Info("Request processed", "duration", time.Since(start), "method", event.HTTPMethod, "path", event.Path)
	}()

	// Get and validate simulation input
	simInput := strings.TrimSpace(event.Body)
	if simInput == "" {
		logger.Warn("Empty request body")
		return transport.NewErrorResponse(http.StatusBadRequest, "Request body cannot be empty")
	}

	logger.Info("Starting simulation", "input_size", len(simInput))
	simulator := sim.NewSimulator()
	result, err := simulator.Run(ctx, simInput)
	if err != nil {
		logger.Error("Simulation failed", "error", err)
		return transport.NewErrorResponse(http.StatusInternalServerError, fmt.Sprintf("Simulation failed: %s", err.Error()))
	}

	response := models.SimulationResponse{Result: result}
	body, err := json.Marshal(response)
	if err != nil {
		logger.Error("Failed to marshal response", "error", err)
		return transport.NewErrorResponse(http.StatusInternalServerError, "Failed to format response")
	}

	return transport.NewResponse(body)
}

func main() {
	lambda.Start(handler)
}
