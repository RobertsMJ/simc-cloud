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
	mux.Handle("GET /health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}))

	return mux
}
