package simc

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/RobertsMJ/simc-cloud/test-utils"
	"github.com/stretchr/testify/suite"
)

type SimcTestSuite struct {
	suite.Suite
}

func TestSimcTestSuite(t *testing.T) {
	suite.Run(t, new(SimcTestSuite))
}

func (s *SimcTestSuite) TestNewSimulationDocument() {

	input, err := os.ReadFile("internal/fixtures/test_dh.simc")
	s.NoError(err)

	res, err := NewSimulationDocument(bytes.NewReader(input))
	s.NoError(err)
	s.NotNil(res)

	path, err := filepath.Abs("internal/fixtures/test_dh_sim_doc.json")
	s.NoError(err)
	expected := test.GoldenValue(s.T(), path, res)
	s.Equal(res, expected)

}
