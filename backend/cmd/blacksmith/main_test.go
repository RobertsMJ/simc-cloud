package main_test

import (
	"testing"

	main "github.com/RobertsMJ/simc-cloud/backend/cmd/blacksmith"
	"github.com/RobertsMJ/simc-cloud/backend/platform"
	"github.com/stretchr/testify/suite"
	"go.uber.org/fx"
)

type GenerateGearsetTestSuite struct {
	suite.Suite
}

func TestGenerateGearsetTestSuite(t *testing.T) {
	suite.Run(t, new(GenerateGearsetTestSuite))
}

func (s *GenerateGearsetTestSuite) TestApp() {
	combinedOptions := append(platform.DefaultOptions(), main.BuildOptions()...)
	s.NoError(fx.ValidateApp(combinedOptions...))
}
