package job

import (
	"context"

	"github.com/RobertsMJ/simc-cloud/backend/models"
)

type Repo interface {
	Create(ctx context.Context, job models.Job) error
}
