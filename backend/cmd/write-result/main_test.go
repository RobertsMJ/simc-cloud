package main_test

import (
	"testing"

	main "github.com/RobertsMJ/simc-cloud/backend/cmd/write-result"
	"github.com/RobertsMJ/simc-cloud/backend/platform"
	"github.com/stretchr/testify/suite"
	"go.uber.org/fx"
)

type WriteResultTestSuite struct {
	suite.Suite
}

func TestWriteResultTestSuite(t *testing.T) {
	suite.Run(t, new(WriteResultTestSuite))
}

func (s *WriteResultTestSuite) TestApp() {
	combinedOptions := append(platform.DefaultOptions(), main.BuildOptions()...)
	s.NoError(fx.ValidateApp(combinedOptions...))
}
