package sim

import (
	"context"

	"github.com/RobertsMJ/simc-cloud/backend/simc"
)

type Input struct {
	Profile string
}

type Simulator interface {
	Run(ctx context.Context, in Input) (simc.Report, error)
}
