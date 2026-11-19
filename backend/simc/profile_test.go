package simc

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/RobertsMJ/simc-cloud/backend/test"
	"github.com/stretchr/testify/suite"
)

type ProfileTestSuite struct {
	suite.Suite
}

func TestProfileTestSuite(t *testing.T) {
	suite.Run(t, new(ProfileTestSuite))
}

func (s *ProfileTestSuite) TestNewProfile_DH() {

	input, err := os.ReadFile("internal/testdata/test_dh_in.simc")
	s.NoError(err)

	res, err := NewProfile(bytes.NewReader(input))
	s.NoError(err)
	s.NotNil(res)

	path, err := filepath.Abs("internal/testdata/test_dh_out.json")
	s.NoError(err)
	expected := test.GoldenValueJSON(s.T(), path, res)
	s.Equal(expected, res)

}

func (s *ProfileTestSuite) TestNewProfile_Warrior() {

	input, err := os.ReadFile("internal/testdata/test_warrior_in.simc")
	s.NoError(err)

	res, err := NewProfile(bytes.NewReader(input))
	s.NoError(err)
	s.NotNil(res)

	path, err := filepath.Abs("internal/testdata/test_warrior_out.json")
	s.NoError(err)
	expected := test.GoldenValueJSON(s.T(), path, res)
	s.Equal(expected, res)

}
