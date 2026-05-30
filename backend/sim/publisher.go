package sim

import (
	"context"

	"github.com/RobertsMJ/simc-cloud-backend/models"
)

type ResultPublisher interface {
	Publish(ctx context.Context, msg models.SimResult) error
}
