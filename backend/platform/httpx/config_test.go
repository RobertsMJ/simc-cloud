package httpx_test

import (
	"testing"
	"time"

	"github.com/RobertsMJ/simc-cloud/backend/platform/httpx"
	"github.com/stretchr/testify/suite"
)

type ConfigTestSuite struct {
	suite.Suite
}

func TestConfigTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}

func (s *ConfigTestSuite) TestDefaults() {
	cfg, err := httpx.NewConfigWithPrefix("admin", "8081")

	s.Require().NoError(err)
	s.Equal(httpx.Config{ID: "admin", Host: "0.0.0.0", Port: "8081", Timeout: 10 * time.Second}, *cfg)
}

func (s *ConfigTestSuite) TestEnvOverrides() {
	s.T().Setenv("ADMIN_PORT", "9000")
	s.T().Setenv("ADMIN_TIMEOUT", "2s")

	cfg, err := httpx.NewConfigWithPrefix("admin", "8081")

	s.Require().NoError(err)
	s.Equal("9000", cfg.Port)
	s.Equal(2*time.Second, cfg.Timeout)
}

func (s *ConfigTestSuite) TestMalformedTimeoutFails() {
	s.T().Setenv("ADMIN_TIMEOUT", "soon")

	_, err := httpx.NewConfigWithPrefix("admin", "8081")

	s.Error(err)
}
