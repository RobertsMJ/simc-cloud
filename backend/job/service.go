package job

import (
	"context"
	"log/slog"

	"github.com/RobertsMJ/simc-cloud/backend/models"
	"github.com/RobertsMJ/simc-cloud/backend/transport/queue"
)

type Service struct {
	log       *slog.Logger
	repo      Repo
	publisher queue.Publisher[models.SimRequest]
}

func NewService(log *slog.Logger, repo Repo, publisher queue.Publisher[models.SimRequest]) *Service {
	return &Service{
		log:       log,
		repo:      repo,
		publisher: publisher,
	}
}

func (s *Service) CreateJob(context.Context, models.Job) (models.Job, error) {
	return models.Job{}, nil
}
