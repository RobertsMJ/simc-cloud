package main_test

import (
	"testing"

	main "github.com/RobertsMJ/simc-cloud/backend/cmd/api"
	"github.com/RobertsMJ/simc-cloud/backend/platform"
	"github.com/stretchr/testify/suite"
	"go.uber.org/fx"
)

type ApiTestSuite struct {
	suite.Suite
}

func TestApiTestSuite(t *testing.T) {
	suite.Run(t, new(ApiTestSuite))
}

func (s *ApiTestSuite) TestApp() {
	combinedOptions := append(platform.DefaultOptions(), main.BuildOptions()...)
	s.NoError(fx.ValidateApp(combinedOptions...))
}
