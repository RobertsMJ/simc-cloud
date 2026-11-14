package job

import (
	"context"
	"log/slog"
	"time"

	"github.com/RobertsMJ/simc-cloud/backend/models"
	"github.com/RobertsMJ/simc-cloud/backend/simc"
	"github.com/RobertsMJ/simc-cloud/backend/transport/queue"
	"github.com/google/uuid"
)

type Service struct {
	log       *slog.Logger
	repo      Repo
	publisher queue.Publisher[models.GenerateGearsetRequest]
}

func NewService(log *slog.Logger, repo Repo, publisher queue.Publisher[models.GenerateGearsetRequest]) *Service {
	return &Service{
		log:       log,
		repo:      repo,
		publisher: publisher,
	}
}

func (s *Service) CreateJob(ctx context.Context, character simc.Character, baseline simc.Loadout, options simc.LoadoutOptions, config simc.Config) (models.Job, error) {

	// Create the job instance that we're working with
	job := models.Job{
		ID:        uuid.NewString(),
		Status:    models.StatusPending,
		CreatedAt: time.Now().Format(time.RFC3339),
	}

	// Create the generate gearset request that will be published to the queue
	gen := models.GenerateGearsetRequest{
		JobID:     job.ID,
		Character: character,
		Baseline:  baseline,
		Options:   options,
		Config:    config,
	}

	if err := s.repo.Create(ctx, job); err != nil {
		return models.Job{}, err
	}

	if err := s.publisher.Publish(ctx, gen); err != nil {
		return models.Job{}, err
	}
	return job, nil
}
