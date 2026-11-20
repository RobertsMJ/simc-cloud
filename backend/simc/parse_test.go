package simc

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/RobertsMJ/simc-cloud/backend/test"
	"github.com/stretchr/testify/suite"
)

type ParseTestSuite struct {
	suite.Suite
}

func TestParseTestSuite(t *testing.T) {
	suite.Run(t, new(ParseTestSuite))
}

func (s *ParseTestSuite) TestNewProfile_DH() {

	input, err := os.ReadFile("testdata/test_dh_in.simc")
	s.NoError(err)

	res, err := Parse(bytes.NewReader(input))
	s.NoError(err)
	s.NotNil(res)

	path, err := filepath.Abs("testdata/test_dh_out.json")
	s.NoError(err)
	expected := test.GoldenValueJSON(s.T(), path, res)
	s.Equal(expected, res)

}

func (s *ParseTestSuite) TestNewProfile_Warrior() {

	input, err := os.ReadFile("testdata/test_warrior_in.simc")
	s.NoError(err)

	res, err := Parse(bytes.NewReader(input))
	s.NoError(err)
	s.NotNil(res)

	path, err := filepath.Abs("testdata/test_warrior_out.json")
	s.NoError(err)
	expected := test.GoldenValueJSON(s.T(), path, res)
	s.Equal(expected, res)

}
