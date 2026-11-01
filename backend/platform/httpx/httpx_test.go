package httpx_test

import (
	"context"
	"log/slog"
	"net/http"
	"testing"

	"github.com/RobertsMJ/simc-cloud/backend/platform/httpx"
	"github.com/stretchr/testify/suite"
)

type HttpxTestSuite struct {
	suite.Suite
}

func TestHttpxTestSuite(t *testing.T) {
	suite.Run(t, new(HttpxTestSuite))
}

type stubHandler struct{}

func (stubHandler) ServeHTTP(http.ResponseWriter, *http.Request) {}

func (s *HttpxTestSuite) TestShutdown() {
	server := httpx.NewServer(
		&httpx.Config{},
		slog.Default(),
		stubHandler{},
	)

	s.Require().NoError(server.Start())
	s.NoError(server.Stop(context.Background()))
	s.NoError(server.Stop(context.Background()))
}
