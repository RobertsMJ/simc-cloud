package sim

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/RobertsMJ/simc-cloud-backend/logger"
)

type Simulator interface {
	Run(ctx context.Context, input string) (string, error)
}

type simulator struct{}

func NewSimulator() Simulator {
	return &simulator{}
}

func (s *simulator) Run(ctx context.Context, input string) (string, error) {
	// Prepare simc command arguments
	args := parseSimcArgs(input)

	// Write args to a temp file
	argsFile, err := os.CreateTemp("", "simc-args.txt")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(argsFile.Name())

	_, err = argsFile.WriteString(strings.Join(args, "\n"))
	if err != nil {
		return "", fmt.Errorf("failed to write simulation args to temp file: %w", err)
	}

	// Execute simc command with context for timeout handling
	logger.Info("Running simulation", "args", argsFile.Name())
	cmd := exec.CommandContext(ctx, "/app/simc", argsFile.Name())

	// Capture both stdout and stderr
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("simc execution failed: %w, output: %s", err, string(output))
	}

	// Return the simulation output
	return strings.TrimSpace(string(output)), nil

}

func parseSimcArgs(input string) []string {
	// Handle windows-encoded newlines
	simcFile := strings.ReplaceAll(input, "\r\n", "\n")
	args := strings.Split(simcFile, "\n")
	// Filter out empty lines and comments
	args = filter(args, func(s string) bool {
		return s != "" && s != "\n" && !strings.HasPrefix(s, "#")
	})
	return args
}

func filter[T any](slice []T, predicate func(T) bool) (filtered []T) {
	for _, item := range slice {
		if predicate(item) {
			filtered = append(filtered, item)
		}
	}
	return
}
