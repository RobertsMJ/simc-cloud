package profile

import (
	"context"
	"strings"

	"github.com/RobertsMJ/simc-cloud/backend/simc"
)

type Service interface {
	Parse(ctx context.Context, raw string) (simc.SimulationDocument, error)
}

type service struct {
}

func NewService() Service {
	return &service{}
}

func (s *service) Parse(ctx context.Context, raw string) (simc.SimulationDocument, error) {
	return simc.NewSimulationDocument(strings.NewReader(raw))
}
