package sim

import (
	"context"

	"github.com/RobertsMJ/simc-cloud-backend/models"
)

type Reader interface {
	GetResult(ctx context.Context, jobID string, gearsetID string) (models.SimResult, error)
	GetResultsByJobID(ctx context.Context, jobID string) ([]models.SimResult, error)
}
