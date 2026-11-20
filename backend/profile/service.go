package profile

import (
	"context"
	"log/slog"
	"strings"

	"github.com/RobertsMJ/simc-cloud/backend/simc"
)

type Service struct {
	log *slog.Logger
}

func NewService(log *slog.Logger) *Service {
	return &Service{
		log: log,
	}
}

func (s Service) Parse(ctx context.Context, raw string) (simc.AddonExport, error) {
	return simc.ParseAddonExport(strings.NewReader(raw))
}
