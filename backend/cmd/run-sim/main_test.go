package main_test

import (
	"testing"

	main "github.com/RobertsMJ/simc-cloud/backend/cmd/run-sim"
	"github.com/RobertsMJ/simc-cloud/backend/platform"
	"github.com/stretchr/testify/suite"
	"go.uber.org/fx"
)

type RunSimTestSuite struct {
	suite.Suite
}

func TestRunSimTestSuite(t *testing.T) {
	suite.Run(t, new(RunSimTestSuite))
}

func (s *RunSimTestSuite) TestApp() {
	combinedOptions := append(platform.DefaultOptions(), main.BuildOptions()...)
	s.NoError(fx.ValidateApp(combinedOptions...))
}
