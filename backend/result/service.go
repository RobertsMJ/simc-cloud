package result

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

var _ queue.Handler[models.SimResult] = (*Service)(nil)

func (s *Service) Handle(ctx context.Context, result models.SimResult) error {
	slog.Info("handling sim result", "result", result)
	return nil
}
