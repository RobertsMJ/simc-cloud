package platform

import "net/http"

type ReadyCheckHandler struct{}

func NewReadyCheckHandler() *ReadyCheckHandler {
	return &ReadyCheckHandler{}
}

func (h *ReadyCheckHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("READY"))
}

func (h *ReadyCheckHandler) Path() string {
	return "/ready"
}
