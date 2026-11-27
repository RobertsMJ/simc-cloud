package result

import (
	"context"

	"github.com/RobertsMJ/simc-cloud/backend/models"
)

type Repo interface {
	RecordResult(ctx context.Context, result models.JobResult) (bool, error)
}
