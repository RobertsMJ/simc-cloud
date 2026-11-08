package api

import (
	"context"
	"log/slog"

	"github.com/RobertsMJ/simc-cloud/backend/profile"
	"github.com/RobertsMJ/simc-cloud/backend/simc"
	"github.com/danielgtaylor/huma/v2"
)

type ProfileRequest struct {
	RawBody []byte
}

type ProfileResponseBody struct {
	Profile simc.SimulationDocument `json:"profile" doc:"Parsed simulation profile"`
}

type ProfileHandler struct {
	log *slog.Logger
	svc profile.Service
}

func NewProfileHandler(svc profile.Service, log *slog.Logger) *ProfileHandler {
	return &ProfileHandler{
		log: log,
		svc: svc,
	}
}

func (h *ProfileHandler) HandleProfileParse(ctx context.Context, request *ProfileRequest) (*Response[ProfileResponseBody], error) {
	doc, err := h.svc.Parse(ctx, string(request.RawBody))
	if err != nil {
		return nil, huma.Error400BadRequest("document parse", err)
	}

	return &Response[ProfileResponseBody]{
		Body: ProfileResponseBody{
			Profile: doc,
		},
	}, nil
}
