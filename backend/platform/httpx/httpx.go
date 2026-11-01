package httpx

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
)

type Server struct {
	server *http.Server
	logger *slog.Logger
	done   chan error
}

func NewServer(config *Config, logger *slog.Logger, handler http.Handler) *Server {
	return &Server{
		server: &http.Server{
			Addr:              net.JoinHostPort(config.Host, config.Port),
			Handler:           handler,
			ReadHeaderTimeout: config.Timeout,
		},
		logger: logger.With("id", config.ID),
		done:   make(chan error, 1),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.server.Addr)
	if err != nil {
		return err
	}
	s.logger.Info("running HTTP server", "addr", s.server.Addr)

	go func() {
		defer close(s.done)
		if err := s.server.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Error("httpx server error", "err", err)
			s.done <- err
			return
		}
		s.done <- nil
	}()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	if err := s.server.Shutdown(ctx); err != nil {
		s.logger.Error("failed to shut down server", "err", err)
		return err
	}
	return nil
}

func (s *Server) Done() <-chan error {
	return s.done
}
