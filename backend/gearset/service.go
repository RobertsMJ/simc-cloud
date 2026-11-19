package gearset

import (
	"log/slog"

	"github.com/RobertsMJ/simc-cloud/backend/models"
)

type Service struct {
	log *slog.Logger
}

func NewService(log *slog.Logger) *Service {
	return &Service{
		log: log,
	}
}

func (s *Service) CountGearsets(baseline models.Gearset, options *models.GearsetOptions) (int, error) {
	// do math
	// simply multiple number of options for each slot
	// trinkets and rings are n choose 2
	// (item sets) * talents * omniums
	return 1, nil
}

func (s *Service) GenerateGearsets(baseline models.Gearset, options *models.GearsetOptions) ([]models.Gearset, error) {
	// lol get wrecked
	return []models.Gearset{baseline}, nil
}
