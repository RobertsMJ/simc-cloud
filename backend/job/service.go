package job

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/RobertsMJ/simc-cloud/backend/models"
	"github.com/RobertsMJ/simc-cloud/backend/simc"
	"github.com/RobertsMJ/simc-cloud/backend/transport/queue"
	"github.com/google/uuid"
)

type GearsetCounter interface {
	CountGearsets(baseline simc.Loadout, options models.Candidates) int
}

type Service struct {
	log            *slog.Logger
	config         *Config
	repo           Repo
	publisher      queue.Publisher[models.GenerateGearsets]
	gearsetService GearsetCounter
}

func NewService(
	log *slog.Logger,
	config *Config,
	repo Repo,
	publisher queue.Publisher[models.GenerateGearsets],
	gearsetService GearsetCounter,
) *Service {
	return &Service{
		log:            log,
		config:         config,
		repo:           repo,
		publisher:      publisher,
		gearsetService: gearsetService,
	}
}

func (s *Service) Create(ctx context.Context, spec models.JobSpec) (models.Job, error) {
	// Check that the number of possible gearsets is within acceptable limits
	if count := s.gearsetService.CountGearsets(spec.Base, spec.Options); count > s.config.MaxGearsets {
		return models.Job{}, fmt.Errorf("too many possible gearsets: %d", count)
	}

	// Create the job instance that we're working with
	job := models.Job{
		ID:        uuid.NewString(),
		Status:    models.StatusPending,
		Spec:      spec,
		CreatedAt: time.Now(),
	}

	// Create the generate gearset request that will be published to the queue
	gen := models.GenerateGearsets{
		JobID: job.ID,
		Spec:  spec,
	}

	if err := s.repo.Create(ctx, job); err != nil {
		return models.Job{}, err
	}

	if err := s.publisher.Publish(ctx, gen); err != nil {
		// TODO: mark the job as failed in the repository
		return models.Job{}, err
	}
	return job, nil
}
