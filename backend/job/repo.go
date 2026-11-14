package job

import (
	"context"

	"github.com/RobertsMJ/simc-cloud/backend/models"
)

// TODO: figure out interface
type Repo interface {
	Create(context.Context, models.Job) error
}
