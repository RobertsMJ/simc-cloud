package job

import (
	"context"
	"log/slog"

	"github.com/RobertsMJ/simc-cloud/backend/models"
	"github.com/RobertsMJ/simc-cloud/backend/transport/queue"
)

type Service struct {
	repo Repo
}

func NewService(repo Repo) *Service {
	return &Service{
		repo: repo,
	}
}

var _ queue.Publisher[models.SimRequest] = (*Service)(nil)

func (s *Service) Publish(ctx context.Context, sim models.SimRequest) error {
	slog.Info("publishing sim request", "result", sim)
	return nil
}
