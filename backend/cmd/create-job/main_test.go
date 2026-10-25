package main_test

import (
	"testing"

	main "github.com/RobertsMJ/simc-cloud/backend/cmd/create-job"
	"github.com/stretchr/testify/suite"
	"go.uber.org/fx"
)

type CreateJobTestSuite struct {
	suite.Suite
}

func TestCreateJobTestSuite(t *testing.T) {
	suite.Run(t, new(CreateJobTestSuite))
}

func (s *CreateJobTestSuite) TestApp() {
	s.NoError(fx.ValidateApp(main.BuildOptions()...))
}
