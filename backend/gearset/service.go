package gearset

import (
	"context"
	"log/slog"

	"github.com/RobertsMJ/simc-cloud/backend/models"
	"github.com/RobertsMJ/simc-cloud/backend/simc"
	"github.com/RobertsMJ/simc-cloud/backend/transport/queue"
)

type Service struct {
	log       *slog.Logger
	publisher queue.Publisher[models.SimRequest]
}

var _ queue.Handler[models.GenerateGearsetRequest] = (*Service)(nil)

func NewService(log *slog.Logger, publisher queue.Publisher[models.SimRequest]) *Service {
	return &Service{
		log:       log,
		publisher: publisher,
	}
}

func (s *Service) Handle(ctx context.Context, req models.GenerateGearsetRequest) error {
	log := s.log.With("job_id", req.JobID)
	log.Info("handling gearset generation")

	profile := simc.Profile{
		Character: req.Character,
		Baseline:  req.Baseline,
		Config:    req.Config,
	}

	doc, err := profile.MarshalSimc()
	if err != nil {
		return err
	}

	if err = s.publisher.Publish(ctx, models.SimRequest{
		JobID: req.JobID,
		Input: doc,
	}); err != nil {
		return err
	}

	return nil
}
