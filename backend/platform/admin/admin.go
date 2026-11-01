package admin

import (
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
		w.Write([]byte("ok"))
	}))
	// TODO: Derive this from readiness checkers
	mux.Handle("/ready", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}))

	return mux
}
