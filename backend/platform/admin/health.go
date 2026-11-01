package admin

import (
	"log/slog"
	"net/http"
)

type HealthCheckHandler struct{}

func NewHealthCheckHandler() *HealthCheckHandler {
	return &HealthCheckHandler{}
}

func (h *HealthCheckHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("ok")); err != nil {
		slog.Warn("failed to write health response", "error", err)
	}
}

func (h *HealthCheckHandler) Path() string {
	return "/health"
}
