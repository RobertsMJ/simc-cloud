package api

import (
	"context"

	genapi "github.com/RobertsMJ/simc-cloud/backend/api/gen"
	"github.com/RobertsMJ/simc-cloud/backend/profile"
)

type ProfileHandler struct {
	svc *profile.Service
}

func NewProfileHandler(svc *profile.Service) *ProfileHandler {
	return &ProfileHandler{svc: svc}
}

func (h *ProfileHandler) PostApiProfile(ctx context.Context, request genapi.PostApiProfileRequestObject) (genapi.PostApiProfileResponseObject, error) {

	// _, err := h.svc.Parse(ctx, request.Body)
	return genapi.PostApiProfile200JSONResponse{
		"req": request.Body,
	}, nil
}
