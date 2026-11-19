package api

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/RobertsMJ/simc-cloud/backend/models"
	"github.com/danielgtaylor/huma/v2"
	"github.com/go-playground/validator/v10"
)

type JobService interface {
	CreateJob(context.Context, models.SimcJobConfig) (models.Job, error)
}

type JobHandler struct {
	log        *slog.Logger
	jobService JobService
}

var _ Handler = (*JobHandler)(nil)

func NewJobHandler(log *slog.Logger, jobService JobService) *JobHandler {
	return &JobHandler{
		log:        log,
		jobService: jobService,
	}
}

func (h JobHandler) RegisterHandler(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "create-job",
		Method:      http.MethodPut,
		Path:        "/jobs",
	}, h.HandleCreateJob)
}

func (h *JobHandler) HandleCreateJob(ctx context.Context, request *Request[models.CreateJobRequestBody]) (*Response[models.CreateJobResponseBody], error) {
	// A job request should contain a character, baseline gear, and selected gear options
	// validate the request to confirm we have all expected configuration
	validate := validator.New()
	if err := validate.Struct(request.Body); err != nil {
		var validationErrors validator.ValidationErrors
		errors.As(err, &validationErrors)
		return nil, validationErrors
	}

	// pass the character, baseline, and options to the job service to handle the business logic
	job, err := h.jobService.CreateJob(ctx, models.SimcJobConfig{
		Character: request.Body.Character,
		Baseline:  request.Body.Baseline,
		Options: models.GearsetOptions{
			Items:   request.Body.Options.Items,
			Talents: request.Body.Options.Talents,
			Omnium:  request.Body.Options.Omnium,
		},
		Config: request.Body.Config,
	})
	if err != nil {
		return nil, err
	}

	return &Response[models.CreateJobResponseBody]{
		Body: models.CreateJobResponseBody{
			Id:     job.ID,
			Status: job.Status,
		},
	}, nil
}
