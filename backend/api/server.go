package api

import (
	"net/http"

	genapi "github.com/RobertsMJ/simc-cloud/backend/api/gen"
	"github.com/RobertsMJ/simc-cloud/backend/profile"
)

type server struct {
	*ProfileHandler
}

func NewHandler(svc *profile.Service) http.Handler {
	sv := &server{
		ProfileHandler: NewProfileHandler(svc),
	}
	strict := genapi.NewStrictHandler(sv, nil)

	return genapi.HandlerFromMux(strict, http.NewServeMux())
}
