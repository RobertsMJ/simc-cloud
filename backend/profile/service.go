package profile

import (
	"context"
	"strings"

	"github.com/RobertsMJ/simc-cloud/backend/simc"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s Service) Parse(ctx context.Context, raw string) (simc.SimulationDocument, error) {
	return simc.NewSimulationDocument(strings.NewReader(raw))
}
