package gearset

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/RobertsMJ/simc-cloud/backend/models"
	"github.com/RobertsMJ/simc-cloud/backend/simc"
	"github.com/RobertsMJ/simc-cloud/backend/transport/queue"
)

type GearsetGenerator interface {
	GenerateGearsets(baseline simc.Loadout, options models.Candidates) ([]simc.Loadout, error)
}

type Handler struct {
	log       *slog.Logger
	publisher queue.Publisher[models.SimRequest]
	generator GearsetGenerator
}

var _ queue.Handler[models.GenerateGearsetsRequest] = (*Handler)(nil)

func NewHandler(log *slog.Logger, publisher queue.Publisher[models.SimRequest], generator GearsetGenerator) *Handler {
	return &Handler{
		log:       log,
		publisher: publisher,
		generator: generator,
	}
}

func (h *Handler) Handle(ctx context.Context, req models.GenerateGearsetsRequest) error {
	log := h.log.With("job_id", req.JobID)
	log.Info("handling gearset generation")

	gearsets, err := h.generator.GenerateGearsets(req.Spec.Base, req.Spec.Options)
	if err != nil {
		return err
	}
	log.Info("generated gearsets", "count", len(gearsets))
	profileSets := make(map[string]simc.Loadout)
	for i, gs := range gearsets {
		profileSets[fmt.Sprintf("set_%d", i)] = gs
	}

	siminput := simc.SimInput{
		Character:   req.Spec.Character,
		Base:        req.Spec.Base,
		Profilesets: profileSets,
		Config:      req.Spec.Config,
	}

	var b strings.Builder
	err = simc.Render(&b, siminput)
	if err != nil {
		return errors.Join(err, errors.New("failed to marshal simulation profile"))
	}

	return h.publisher.Publish(ctx, models.SimRequest{
		JobID: req.JobID,
		Input: b.String(),
	})
}
