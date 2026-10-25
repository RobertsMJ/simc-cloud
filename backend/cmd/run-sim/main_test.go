package main_test

import (
	"testing"

	main "github.com/RobertsMJ/simc-cloud/backend/cmd/run-sim"
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
	s.NoError(fx.ValidateApp(main.BuildOptions()...))
}
