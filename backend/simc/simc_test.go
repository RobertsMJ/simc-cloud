package simc

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

type SimcTestSuite struct {
	suite.Suite
}

func TestSimcTestSuite(t *testing.T) {
	suite.Run(t, new(SimcTestSuite))
}

func (s *SimcTestSuite) TestNewSimulationDocument() {

	input, err := os.ReadFile("fixtures/test_dh.simc")
	s.NoError(err)

	res, err := NewSimulationDocument(bytes.NewReader(input))
	s.NoError(err)
	s.NotNil(res)

}
