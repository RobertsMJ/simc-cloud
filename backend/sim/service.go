package sim

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"runtime"
	"slices"
	"strings"

	"github.com/RobertsMJ/simc-cloud/backend/models"
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

	slog.Debug(fmt.Sprintf("Starting simulation for JobID: %s, GearsetID: %s", request.JobID, request.GearsetID))

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
	threads := fmt.Sprintf("threads=%d", runtime.GOMAXPROCS(0))
	slog.Debug("running simulation with threads", "threads", runtime.GOMAXPROCS(0))
	args = append(args, "json2="+outputPath, "html=/dev/null", "report_details=0", "target_error=0.05", threads)

	slog.Debug("parsed simc arguments", "args", args)

	cmd := exec.CommandContext(ctx, "/app/simc", args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return models.SimResult{}, fmt.Errorf("failed to get stdout pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		slog.Error("failed to start simc command", "err", err, "stderr", stderr.String())
		return models.SimResult{}, fmt.Errorf("simc execution failed: %w, stderr: %s", err, stderr.String())
	}

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		txt := strings.TrimSpace(scanner.Text())
		if txt != "" {
			slog.Debug(scanner.Text())
		}
	}
	if err := scanner.Err(); err != nil {
		slog.Error("scanner errored", "err", err)
	}

	if err := cmd.Wait(); err != nil {
		if ctxErr := context.Cause(ctx); ctxErr != nil {
			slog.Error("simc command terminated by context", "ctxErr", ctxErr, "err", err)
			return models.SimResult{}, errors.Join(ctxErr, err)
		}
		slog.Error("simc command failed", "err", err, "stderr", stderr.String())
		return models.SimResult{}, err
	}

	outputBytes, err := os.ReadFile(outputPath)
	if err != nil {
		return models.SimResult{}, fmt.Errorf("failed to read simc output: %w", err)
	}

	var envelope simcOutputEnvelope
	if err := json.Unmarshal(outputBytes, &envelope); err != nil {
		return models.SimResult{}, fmt.Errorf("failed to parse simc output: %w", err)
	}

	slog.Info("sim completed successfully", "job_id", request.JobID, "gearset_id", request.GearsetID, "threads", threads, "cpu_time", envelope.Sim.Statistics.ElapsedCpuSeconds)

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
