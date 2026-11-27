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
	resultPub queue.Publisher[models.SimCompleted]
	config    *Config
}

var _ queue.Handler[models.RunSim] = (*Service)(nil)

func NewService(simulator Simulator, resultPub queue.Publisher[models.SimCompleted], config *Config) *Service {
	return &Service{
		simulator: simulator,
		resultPub: resultPub,
		config:    config,
	}
}

func (s *Service) Handle(ctx context.Context, req models.RunSim) error {
	slog.Info("handling sim message", "job_id", req.JobID)
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
	if errors.Is(simCtx.Err(), context.DeadlineExceeded) {
		err = context.DeadlineExceeded
	}

	simResult := models.SimCompleted{
		JobID:    req.JobID,
		Gearsets: req.Gearsets,
	}
	if err != nil {
		simResult.Status = models.StatusError
		errMsg := err.Error()
		simResult.ErrorMessage = &errMsg
		slog.Error("sim error", "job_id", req.JobID, "error", err)
	} else {
		simResult.Status = models.StatusDone
		simResult.Result = out
	}

	publishCtx, publishCancel := context.WithTimeout(ctx, 1*time.Second)
	defer publishCancel()

	return s.resultPub.Publish(publishCtx, simResult)
}
