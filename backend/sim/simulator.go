package sim

import (
	"context"
	"encoding/json"

	"github.com/RobertsMJ/simc-cloud/simc"
)

type Input struct {
	Profile string
}

type Output struct {
	Statistics simc.Statistics
	Report     json.RawMessage
}

type Simulator interface {
	Run(ctx context.Context, in Input) (Output, error)
}
