package api

import (
	"context"

	"github.com/RobertsMJ/simc-cloud/backend/profile"
)

type ProfileRequestBody struct {
	Profile string `json:"profile" doc:"Simulation profile to parse"`
}

type ProfileResponseBody struct {
	Profile map[string]any `json:"profile" doc:"Parsed simulation profile"`
}

type ProfileHandler struct {
	svc *profile.Service
}

func NewProfileHandler(svc *profile.Service) *ProfileHandler {
	return &ProfileHandler{svc: svc}
}

func (h *ProfileHandler) HandleProfileParse(ctx context.Context, request *Request[ProfileRequestBody]) (*Response[ProfileResponseBody], error) {
	resp := ProfileResponseBody{
		Profile: map[string]any{
			"req": request.Body,
		},
	}
	return &Response[ProfileResponseBody]{
		Body: resp,
	}, nil
}
