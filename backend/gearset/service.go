package gearset

import (
	"log/slog"

	"github.com/RobertsMJ/simc-cloud/backend/models"
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

func (s *Service) CountGearsets(baseline simc.Loadout, options models.Candidates) (int, error) {
	// do math
	// simply multiple number of options for each slot
	// trinkets and rings are n choose 2
	// (item sets) * talents * omniums
	return 1, nil
}

// GenerateGearsets does not include the baseline loadout in the returned gearsets
func (s *Service) GenerateGearsets(baseline simc.Loadout, options models.Candidates) ([]simc.Loadout, error) {
	// lol get wrecked
	return []simc.Loadout{}, nil
}
