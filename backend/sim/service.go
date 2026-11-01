package sim

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/RobertsMJ/simc-cloud/backend/models"
	"github.com/RobertsMJ/simc-cloud/backend/transport/queue"
)

type Service struct {
	simulator Simulator
	resultPub queue.Publisher[models.SimResult]
	config    *Config
}

var _ queue.Handler[models.SimRequest] = (*Service)(nil)

func NewService(simulator Simulator, resultPub queue.Publisher[models.SimResult], config *Config) *Service {
	return &Service{
		simulator: simulator,
		resultPub: resultPub,
		config:    config,
	}
}

func (s *Service) Handle(ctx context.Context, req models.SimRequest) error {
	slog.Info("handling sim message", "job_id", req.JobID, "gearset_id", req.GearsetID)
	simCtx, simCancel := context.WithTimeout(ctx, s.config.MaxRuntime)
	defer simCancel()

	out, err := s.simulator.Run(
		simCtx,
		Input{
			Profile: req.Input,
		},
	)
	if ctx.Err() != nil {
		return ctx.Err()
	}
	if simCtx.Err() != nil {
		return simCtx.Err()
	}

	simResult := models.SimResult{
		JobID:     req.JobID,
		GearsetID: req.GearsetID,
		Metadata:  req.Metadata,
	}
	if err != nil {
		simResult.Status = models.StatusError
		errMsg := err.Error()
		simResult.ErrorMessage = &errMsg
		if errors.Is(err, context.DeadlineExceeded) {
			slog.Error("sim deadline exceeded", "job_id", req.JobID, "gearset_id", req.GearsetID)
			return queue.ErrFatal
		}
		slog.Error("sim failed", "job_id", req.JobID, "gearset_id", req.GearsetID, "error", err)
	} else {
		simResult.Status = models.StatusCompleted
		simResult.Statistics = out.Statistics
		simResult.Result = out.Report
	}

	publishCtx, publishCancel := context.WithTimeout(ctx, 1*time.Second)
	defer publishCancel()

	return s.resultPub.Publish(publishCtx, simResult)
}
