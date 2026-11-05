package simcexec

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"strings"

	"github.com/RobertsMJ/simc-cloud/backend/sim"
	"github.com/RobertsMJ/simc-cloud/backend/simc"
)

type SimcExec struct {
	config *Config
}

var _ sim.Simulator = (*SimcExec)(nil)

func New(config *Config) *SimcExec {
	return &SimcExec{
		config: config,
	}
}

func (s SimcExec) Run(ctx context.Context, in sim.Input) (sim.Output, error) {
	if in.Profile == "" {
		return sim.Output{}, fmt.Errorf("empty simc profile")
	}

	tmpDir, err := os.MkdirTemp("", "simc-*")
	if err != nil {
		return sim.Output{}, fmt.Errorf("create temp dir: %w", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			slog.Warn("failed to remove simc temp dir", "path", tmpDir, "error", err)
		}
	}()

	inputPath := tmpDir + "/input.simc"
	if err := os.WriteFile(inputPath, []byte(in.Profile), 0644); err != nil {
		return sim.Output{}, fmt.Errorf("write input file: %w", err)
	}

	slog.Debug("running simulation with threads", "threads", s.config.Threads)

	outputPath := tmpDir + "/output.json"
	cmd := exec.CommandContext(ctx, s.config.BinPath, inputPath, "json2="+outputPath, "html=/dev/null", "report_details=0", fmt.Sprintf("threads=%d", s.config.Threads))
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return sim.Output{}, fmt.Errorf("stdout pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return sim.Output{}, fmt.Errorf("simc execute: %w", err)
	}

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		if line := strings.TrimSpace(scanner.Text()); line != "" {
			slog.Debug(line)
		}
	}
	if err := scanner.Err(); err != nil {
		slog.Warn("stopped logging simc output", "err", err)
	}
	// Drain whatever the scanner left unread so simc can't block on a full pipe.
	// Must happen before Wait, which closes the read end of the pipe.
	if _, err := io.Copy(io.Discard, stdout); err != nil {
		slog.Debug("drain simc stdout", "err", err)
	}

	if err := cmd.Wait(); err != nil {
		if ctxErr := context.Cause(ctx); ctxErr != nil {
			return sim.Output{}, fmt.Errorf("simc cancelled: %w", ctxErr)
		}

		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return sim.Output{}, &ExecError{ExitCode: exitErr.ExitCode(), Stderr: stderr.String(), Err: err}
		}
		return sim.Output{}, fmt.Errorf("wait for simc: %w", err)
	}

	outputBytes, err := os.ReadFile(outputPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return sim.Output{}, fmt.Errorf("simc produced no report: %s", strings.TrimSpace(stderr.String()))
		}
		return sim.Output{}, fmt.Errorf("read simc output: %w", err)
	}

	res, err := simc.ParseReport(bytes.NewReader(outputBytes))
	if err != nil {
		return sim.Output{}, fmt.Errorf("report parse: %w", err)
	}

	slog.Info("sim completed successfully", "threads", s.config.Threads, "cpu_time", res.Sim.Statistics.ElapsedCpuSeconds)

	return sim.Output{
		Statistics: res.Sim.Statistics,
		Report:     outputBytes,
	}, nil
}

type ExecError struct {
	ExitCode int
	Stderr   string
	Err      error
}

func (e *ExecError) Error() string {
	msg := strings.TrimSpace(e.Stderr)
	if msg == "" {
		return fmt.Sprintf("simc exited with code %d", e.ExitCode)
	}
	return fmt.Sprintf("simc exited with code %d: %s", e.ExitCode, msg)
}

func (e *ExecError) Unwrap() error { return e.Err }
