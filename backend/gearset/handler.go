package gearset

import (
	"context"
	"errors"
	"log/slog"

	"github.com/RobertsMJ/simc-cloud/backend/models"
	"github.com/RobertsMJ/simc-cloud/backend/simc"
	"github.com/RobertsMJ/simc-cloud/backend/transport/queue"
)

type GearsetGenerator interface {
	GenerateGearsets(baseline models.Gearset, options *models.GearsetOptions) ([]models.Gearset, error)
}

type Handler struct {
	log       *slog.Logger
	publisher queue.Publisher[models.SimRequest]
	generator GearsetGenerator
}

var _ queue.Handler[models.GenerateGearsetRequest] = (*Handler)(nil)

func NewHandler(log *slog.Logger, publisher queue.Publisher[models.SimRequest], generator GearsetGenerator) *Handler {
	return &Handler{
		log:       log,
		publisher: publisher,
		generator: generator,
	}
}

func (h *Handler) Handle(ctx context.Context, req models.GenerateGearsetRequest) error {
	log := h.log.With("job_id", req.JobID)
	log.Info("handling gearset generation")

	gearsets, err := h.generator.GenerateGearsets(req.Baseline, &req.Options)
	if err != nil {
		return err
	}
	log.Info("generated gearsets", "count", len(gearsets))

	simcLoadouts := make([]simc.Loadout, 0, len(gearsets))
	for _, gs := range gearsets {
		simcLoadouts = append(simcLoadouts, gs.ToSimcLoadout())
	}

	profile := simc.SimulationProfile{
		Character: req.Character,
		Baseline:  req.Baseline.ToSimcLoadout(),
		Options:   simcLoadouts,
		Config:    req.Config,
	}

	doc, err := profile.MarshalSimc()
	if err != nil {
		return errors.Join(err, errors.New("failed to marshal simulation profile"))
	}

	if err = h.publisher.Publish(ctx, models.SimRequest{
		JobID: req.JobID,
		Input: doc,
	}); err != nil {
		return err
	}

	return nil
}
