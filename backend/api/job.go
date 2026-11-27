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
	Create(context.Context, models.JobSpec) (models.Job, error)
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

type CreateJobRequest = Request[models.JobSpec]
type CreateJobResponseBody struct {
	Id string `json:"id"`
}
type CreateJobResponse = Response[CreateJobResponseBody]

func (h *JobHandler) HandleCreateJob(ctx context.Context, request *CreateJobRequest) (*CreateJobResponse, error) {
	// A job request should contain a character, baseline gear, and selected gear options
	// validate the request to confirm we have all expected configuration
	validate := validator.New()
	if err := validate.Struct(request.Body); err != nil {
		var validationErrors validator.ValidationErrors
		errors.As(err, &validationErrors)
		return nil, validationErrors
	}

	// pass the character, baseline, and options to the job service to handle the business logic
	job, err := h.jobService.Create(ctx, request.Body)
	if err != nil {
		return nil, err
	}

	return &CreateJobResponse{
		Body: CreateJobResponseBody{
			Id: job.ID,
		},
	}, nil
}
