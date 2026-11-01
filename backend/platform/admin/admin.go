package admin

import (
	"log/slog"
	"net/http"

	"github.com/RobertsMJ/simc-cloud/backend/platform/httpx"
)

func NewAdminServerConfig() (*httpx.Config, error) {
	return httpx.NewConfigWithPrefix("admin", "8081")
}

func NewAdminHandler() http.Handler {
	mux := http.NewServeMux()

	// TODO: Derive this from health checkers
	mux.Handle("/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("ok")); err != nil {
			slog.Warn("failed to write health response", "error", err)
		}
	}))
	// TODO: Derive this from readiness checkers
	mux.Handle("/ready", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("ok")); err != nil {
			slog.Warn("failed to write readiness response", "error", err)
		}
	}))

	return mux
}
