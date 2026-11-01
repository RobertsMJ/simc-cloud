package admin

import (
	"log/slog"
	"net/http"
)

type ReadyCheckHandler struct{}

func NewReadyCheckHandler() *ReadyCheckHandler {
	return &ReadyCheckHandler{}
}

func (h *ReadyCheckHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("ok")); err != nil {
		slog.Warn("failed to write readiness response", "error", err)
	}
}

func (h *ReadyCheckHandler) Path() string {
	return "/ready"
}
