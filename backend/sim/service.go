package sim

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/RobertsMJ/simc-cloud-backend/logger"
	"github.com/RobertsMJ/simc-cloud-backend/models"
)

type SimInput = models.SimulationRequest
type SimOutput = models.SimulationResponse

type Simulator interface {
	Run(ctx context.Context, input *SimInput) (SimOutput, error)
}

type simulator struct{}

func NewSimulator() Simulator {
	return &simulator{}
}

func (s *simulator) Run(ctx context.Context, request *SimInput) (SimOutput, error) {
	if request == nil {
		return models.SimulationResponse{}, fmt.Errorf("input cannot be nil")
	}
	// Prepare simc command arguments
	args, err := parseSimcArgs(&request.Input)
	if err != nil {
		return models.SimulationResponse{}, fmt.Errorf("failed to parse simc arguments: %w", err)
	}

	args = append(args, "json2=stdout")
	args = append(args, "report_details=0")

	// Write args to a temp file
	// argsFile, err := os.CreateTemp("", "simc-args.txt")
	// if err != nil {
	// 	return models.SimulationResponse{}, fmt.Errorf("failed to create temp file: %w", err)
	// }
	// defer os.Remove(argsFile.Name())

	// _, err = argsFile.WriteString(strings.Join(args, "\n"))
	// if err != nil {
	// 	return models.SimulationResponse{}, fmt.Errorf("failed to write simulation args to temp file: %w", err)
	// }

	// Execute simc command with context for timeout handling
	// logger.Info("Running simulation", "args", argsFile.Name())
	cmd := exec.CommandContext(ctx, "/app/simc", args...)
	// cmd := exec.CommandContext(ctx, "/app/simc", argsFile.Name())

	// Capture both stdout and stderr
	output, err := cmd.CombinedOutput()
	if err != nil {
		return models.SimulationResponse{}, fmt.Errorf("simc execution failed: %w, output: %s", err, string(output))
	}
	logger.Info("sim results: " + string(output))

	// Return the simulation output
	return models.SimulationResponse{
		RequestID: request.RequestID,
		GearsetID: request.GearsetID,
		Metadata:  request.Metadata,
		Result:    string(output),
	}, nil
}

func parseSimcArgs(input *string) ([]string, error) {
	if input == nil {
		return nil, fmt.Errorf("simc input string cannot be nil")
	}
	// Handle windows-encoded newlines
	simcFile := strings.ReplaceAll(*input, "\r\n", "\n")
	args := strings.Split(simcFile, "\n")
	// Filter out empty lines and comments
	args = filter(args, func(s string) bool {
		return s != "" && s != "\n" && !strings.HasPrefix(s, "#")
	})
	return args, nil
}

func filter[T any](slice []T, predicate func(T) bool) (filtered []T) {
	for _, item := range slice {
		if predicate(item) {
			filtered = append(filtered, item)
		}
	}
	return
}
