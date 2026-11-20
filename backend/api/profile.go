package api

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/RobertsMJ/simc-cloud/backend/simc"
	"github.com/danielgtaylor/huma/v2"
)

type ProfileRequest struct {
	RawBody []byte
}

type ProfileResponseBody struct {
	Profile simc.AddonExport `json:"profile" doc:"Parsed simc addon export"`
}

type ProfileParser interface {
	Parse(ctx context.Context, raw string) (simc.AddonExport, error)
}

type ProfileHandler struct {
	log *slog.Logger
	svc ProfileParser
}

var _ Handler = (*ProfileHandler)(nil)

func NewProfileHandler(log *slog.Logger, svc ProfileParser) *ProfileHandler {
	return &ProfileHandler{
		log: log,
		svc: svc,
	}
}

func (h *ProfileHandler) RegisterHandler(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "parse-profile",
		Method:      http.MethodPost,
		Path:        "/profile",
	}, h.HandleProfileParse)
}

func (h *ProfileHandler) HandleProfileParse(ctx context.Context, request *ProfileRequest) (*Response[ProfileResponseBody], error) {
	doc, err := h.svc.Parse(ctx, string(request.RawBody))
	if err != nil {
		if errors.Is(err, simc.ParseError{}) {
			return nil, huma.Error422UnprocessableEntity("failed to parse profile", err)
		}
		return nil, huma.Error500InternalServerError("document parse", err)
	}

	return &Response[ProfileResponseBody]{
		Body: ProfileResponseBody{
			Profile: doc,
		},
	}, nil
}
