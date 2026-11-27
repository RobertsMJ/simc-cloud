package result

import (
	"context"
	"log/slog"

	"github.com/RobertsMJ/simc-cloud/backend/models"
	"github.com/RobertsMJ/simc-cloud/backend/transport/queue"
)

type Service struct {
	repo Repo
}

func NewService(repo Repo) *Service {
	return &Service{
		repo: repo,
	}
}

var _ queue.Handler[models.SimCompleted] = (*Service)(nil)

func (s *Service) Handle(ctx context.Context, result models.SimCompleted) error {
	slog.Info("handling sim result", "result", result)

	gs := make([]models.GearsetResult, 0, len(result.Result.Sim.Profilesets.Results))
	for _, res := range result.Result.Sim.Profilesets.Results {
		gs = append(gs, models.GearsetResult{
			JobID:     result.JobID,
			ID:        res.Name,
			Loadout:   result.Gearsets[res.Name],
			DPS:       res.Mean,
			DPSStdErr: res.MeanStdDev,
		})
	}

	_, err := s.repo.RecordResult(ctx, models.JobResult{
		JobID:             result.JobID,
		SimcVersion:       result.Result.GitRevision,
		BaselineDPS:       result.Result.Sim.Statistics.RaidDPS.Mean,
		BaselineDPSStdErr: result.Result.Sim.Statistics.RaidDPS.StdDev,
		Gearsets:          gs,
	})
	return err
}
