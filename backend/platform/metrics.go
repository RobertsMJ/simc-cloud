package platform

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type MetricsHandler struct {
	Registry *prometheus.Registry
}

func NewMetricsHandler() *MetricsHandler {
	reg := prometheus.NewRegistry()
	reg.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)
	return &MetricsHandler{
		Registry: reg,
	}
}

func (h *MetricsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	promhttp.HandlerFor(h.Registry, promhttp.HandlerOpts{}).ServeHTTP(w, r)
}

func (h *MetricsHandler) Path() string {
	return "/metrics"
}
