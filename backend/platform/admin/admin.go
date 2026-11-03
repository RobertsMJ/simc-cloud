package admin

import (
	"log/slog"
	"net/http"

	"github.com/RobertsMJ/simc-cloud/backend/platform/httpx"
	"github.com/RobertsMJ/simc-cloud/backend/platform/readiness"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewAdminServerConfig() (*httpx.Config, error) {
	return httpx.NewConfigWithPrefix("admin", "8081")
}

func NewAdminHandler(hc *readiness.Handler, metrics prometheus.Gatherer) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("ok")); err != nil {
			slog.Warn("failed to write health response", "error", err)
		}
	}))
	mux.Handle("/ready", hc)
	mux.Handle("/metrics", promhttp.HandlerFor(metrics, promhttp.HandlerOpts{}))

	return mux
}
