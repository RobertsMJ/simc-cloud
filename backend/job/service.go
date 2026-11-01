package job

import (
	"github.com/RobertsMJ/simc-cloud/backend/models"
	"github.com/RobertsMJ/simc-cloud/backend/transport/queue"
)

type Service struct {
	repo      Repo
	publisher queue.Publisher[models.SimRequest]
}

func NewService(repo Repo, publisher queue.Publisher[models.SimRequest]) *Service {
	return &Service{
		repo:      repo,
		publisher: publisher,
	}
}
