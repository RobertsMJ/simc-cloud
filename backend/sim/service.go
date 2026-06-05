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

	"github.com/RobertsMJ/simc-cloud-backend/models"
)

type simcOutputEnvelope struct {
	Sim struct {
		Statistics models.SimStatistics `json:"statistics"`
	} `json:"sim"`
}

type Simulator interface {
	Run(ctx context.Context, input *models.SimRequest) (models.SimResult, error)
}

type simulator struct {
}

func NewSimulator() Simulator {
	return &simulator{}
}

func (s *simulator) Run(ctx context.Context, request *models.SimRequest) (models.SimResult, error) {
	if request == nil {
		return models.SimResult{}, fmt.Errorf("input cannot be nil")
	}

	args, err := parseSimcArgs(&request.Input)
	if err != nil {
		return models.SimResult{}, fmt.Errorf("failed to parse simc arguments: %w", err)
	}

	tmpDir, err := os.MkdirTemp("", "simc-*")
	if err != nil {
		return models.SimResult{}, fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	outputPath := tmpDir + "/output.json"
	args = append(args, "json2="+outputPath, "html=/dev/null", "report_details=0")

	cmd := exec.CommandContext(ctx, "/app/simc", args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return models.SimResult{}, fmt.Errorf("simc execution failed: %w, stderr: %s", err, stderr.String())
	}

	outputBytes, err := os.ReadFile(outputPath)
	if err != nil {
		return models.SimResult{}, fmt.Errorf("failed to read simc output: %w", err)
	}

	var envelope simcOutputEnvelope
	if err := json.Unmarshal(outputBytes, &envelope); err != nil {
		return models.SimResult{}, fmt.Errorf("failed to parse simc output: %w", err)
	}

	return models.SimResult{
		JobID:      request.JobID,
		GearsetID:  request.GearsetID,
		Status:     models.StatusCompleted,
		Statistics: envelope.Sim.Statistics,
		Metadata:   request.Metadata,
		Result:     outputBytes,
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
