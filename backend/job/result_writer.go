package job

import (
	"context"

	"github.com/RobertsMJ/simc-cloud-backend/models"
)

type ResultWriter interface {
	WriteResult(ctx context.Context, result models.SimResult) error
}
