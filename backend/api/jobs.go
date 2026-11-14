package api

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/RobertsMJ/simc-cloud/backend/models"
	"github.com/RobertsMJ/simc-cloud/backend/simc"
	"github.com/danielgtaylor/huma/v2"
)

type CreateJobRequestBody struct {
	Character simc.Character      `json:"character"`
	Baseline  simc.Loadout        `json:"baseline"`
	Options   simc.LoadoutOptions `json:"options"`
	Config    map[string]string   `json:"config"`
}

type CreateJobResponseBody struct {
	Id string `json:"id"`
}

type JobService interface {
	CreateJob(context.Context, models.Job) (models.Job, error)
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
		Path:        "/job",
	}, h.HandleCreateJob)
}

func (h *JobHandler) HandleCreateJob(ctx context.Context, request *Request[CreateJobRequestBody]) (*Response[CreateJobResponseBody], error) {

	return &Response[CreateJobResponseBody]{
		Body: CreateJobResponseBody{
			Id: "created-job-id",
		},
	}, nil
}
