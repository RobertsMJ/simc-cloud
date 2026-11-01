package config_test

import (
	"testing"
	"time"

	"github.com/RobertsMJ/simc-cloud/backend/config"
	"github.com/stretchr/testify/suite"
)

type ConfigTestSuite struct {
	suite.Suite
}

func TestConfigTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}

type sample struct {
	Name    string        `env:"NAME,required"`
	Count   int           `env:"COUNT" envDefault:"3"`
	Timeout time.Duration `env:"TIMEOUT" envDefault:"5s"`
}

func (s *ConfigTestSuite) TestParse_MissingRequiredFails() {
	_, err := config.Parse[sample]("T_")
	s.Error(err)
}

func (s *ConfigTestSuite) TestParse_AppliesPrefixAndDefaults() {
	s.T().Setenv("T_NAME", "x")

	cfg, err := config.Parse[sample]("T_")

	s.Require().NoError(err)
	s.Equal(sample{Name: "x", Count: 3, Timeout: 5 * time.Second}, *cfg)
}

func (s *ConfigTestSuite) TestParse_MalformedValueFails() {
	s.T().Setenv("T_NAME", "x")
	s.T().Setenv("T_TIMEOUT", "soon")

	_, err := config.Parse[sample]("T_")

	s.Error(err)
}
