package job

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/RobertsMJ/simc-cloud/backend/models"
	"github.com/RobertsMJ/simc-cloud/backend/transport/queue"
	"github.com/google/uuid"
)

type GearsetCounter interface {
	CountGearsets(baseline models.Gearset, options *models.GearsetOptions) (int, error)
}

type Service struct {
	log            *slog.Logger
	config         *Config
	repo           Repo
	publisher      queue.Publisher[models.GenerateGearsetRequest]
	gearsetService GearsetCounter
}

func NewService(log *slog.Logger, config *Config, repo Repo, publisher queue.Publisher[models.GenerateGearsetRequest], gearsetService GearsetCounter) *Service {
	return &Service{
		log:            log,
		config:         config,
		repo:           repo,
		publisher:      publisher,
		gearsetService: gearsetService,
	}
}

func (s *Service) CreateJob(ctx context.Context, jobCfg models.SimcJobConfig) (models.Job, error) {
	// Check that the number of possible gearsets is within acceptable limits
	if count, err := s.gearsetService.CountGearsets(jobCfg.Baseline, &models.GearsetOptions{
		Items:   jobCfg.Options.Items,
		Talents: jobCfg.Options.Talents,
		Omnium:  jobCfg.Options.Omnium,
	}); err != nil {
		return models.Job{}, err
	} else if count > s.config.MaxGearsets {
		return models.Job{}, fmt.Errorf("too many possible gearsets: %d", count)
	}

	// Create the job instance that we're working with
	job := models.Job{
		ID:        uuid.NewString(),
		Status:    models.StatusPending,
		CreatedAt: time.Now().Format(time.RFC3339),
	}

	// Create the generate gearset request that will be published to the queue
	gen := models.GenerateGearsetRequest{
		JobID:     job.ID,
		Character: jobCfg.Character,
		Baseline:  jobCfg.Baseline,
		Config:    jobCfg.Config,
	}

	if err := s.repo.Create(ctx, job); err != nil {
		return models.Job{}, err
	}

	if err := s.publisher.Publish(ctx, gen); err != nil {
		return models.Job{}, err
	}
	return job, nil
}
