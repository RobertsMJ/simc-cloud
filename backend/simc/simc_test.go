package simc

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/RobertsMJ/simc-cloud/backend/test"
	"github.com/stretchr/testify/suite"
)

type SimcTestSuite struct {
	suite.Suite
}

func TestSimcTestSuite(t *testing.T) {
	suite.Run(t, new(SimcTestSuite))
}

func (s *SimcTestSuite) TestMarshal() {
	input, err := os.ReadFile("internal/testdata/test_warrior_sim_in.json")
	s.NoError(err)

	var simProfile SimulationProfile
	err = json.Unmarshal(input, &simProfile)
	s.NoError(err)

	path, err := filepath.Abs("internal/testdata/test_warrior_sim_out.simc")
	s.NoError(err)

	actual, err := simProfile.MarshalSimc()
	s.NoError(err)
	expected := test.GoldenValueRaw(s.T(), path, []byte(actual))
	s.Equal(string(expected), actual)
}

func (s *SimcTestSuite) TestMarshal_Warrior() {
	// input, err := os.ReadFile("internal/testdata/test_warrior_in.simc")
	// s.NoError(err)

	// res, err := NewProfile(bytes.NewReader(input))
	// s.NoError(err)
	// s.NotNil(res)

	// path, err := filepath.Abs("internal/testdata/test_warrior_out.json")
	// s.NoError(err)
	// expected := test.GoldenValueJSON(s.T(), path, res)
	// s.Equal(expected, res)
}
