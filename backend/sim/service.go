package sim

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"slices"
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

	args, err := parseSimcArgs(&request.Input)
	if err != nil {
		return models.SimulationResponse{}, fmt.Errorf("failed to parse simc arguments: %w", err)
	}

	tmpDir, err := os.MkdirTemp("", "simc-*")
	if err != nil {
		return models.SimulationResponse{}, fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	outputPath := tmpDir + "/output.json"
	args = append(args, "json2="+outputPath, "html=/dev/null", "report_details=0")

	cmd := exec.CommandContext(ctx, "/app/simc", args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return models.SimulationResponse{}, fmt.Errorf("simc execution failed: %w, stderr: %s", err, stderr.String())
	}

	outputBytes, err := os.ReadFile(outputPath)
	if err != nil {
		return models.SimulationResponse{}, fmt.Errorf("failed to read simc output: %w", err)
	}

	var result map[string]any
	if err := json.Unmarshal(outputBytes, &result); err != nil {
		return models.SimulationResponse{}, fmt.Errorf("failed to parse simc output: %w", err)
	}

	logger.Debug("simulation complete", "request_id", request.RequestID, "gearset_id", request.GearsetID, "result", result)

	return models.SimulationResponse{
		RequestID: request.RequestID,
		GearsetID: request.GearsetID,
		Metadata:  request.Metadata,
		Result:    result,
	}, nil
}

func parseSimcArgs(input *string) ([]string, error) {
	if input == nil {
		return nil, fmt.Errorf("simc input string cannot be nil")
	}
	args := strings.Split(strings.ReplaceAll(*input, "\r\n", "\n"), "\n")
	return slices.DeleteFunc(args, func(s string) bool {
		return s == "" || strings.HasPrefix(s, "#")
	}), nil
}
