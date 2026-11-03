package readiness

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type Status string

const (
	StatusHealthy   Status = "healthy"
	StatusUnhealthy Status = "unhealthy"
)

type CheckResult struct {
	Name   string
	Status Status
}

type Checker interface {
	Check() CheckResult
}

type Handler struct {
	checkers []Checker
	logger   *slog.Logger
}

func NewHandler(checkers []Checker, logger *slog.Logger) *Handler {
	return &Handler{
		checkers: checkers,
		logger:   logger,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	status := http.StatusOK
	res := map[string]Status{}
	for _, checker := range h.checkers {
		result := checker.Check()
		res[result.Name] = result.Status
		if result.Status == StatusUnhealthy {
			status = http.StatusServiceUnavailable
		}
	}

	w.Header().Add("content-type", "application/json; charset=utf-8")
	body, err := json.Marshal(res)
	if err != nil {
		h.logger.Warn("failed to marshal readiness response", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		w.WriteHeader(status)
	}
	if _, err := w.Write(body); err != nil {
		h.logger.Warn("failed to write readiness response", "error", err)
	}
}
