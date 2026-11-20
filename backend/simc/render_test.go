package simc

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/RobertsMJ/simc-cloud/backend/test"
	"github.com/stretchr/testify/suite"
)

type RenderTestSuite struct {
	suite.Suite
}

func TestRenderTestSuite(t *testing.T) {
	suite.Run(t, new(RenderTestSuite))
}

func (s *RenderTestSuite) TestRender_Warrior() {
	input, err := os.ReadFile("testdata/test_warrior_sim_in.json")
	s.NoError(err)

	var simProfile SimInput
	err = json.Unmarshal(input, &simProfile)
	s.NoError(err)

	path, err := filepath.Abs("testdata/test_warrior_sim_out.simc")
	s.NoError(err)

	var b strings.Builder
	err = Render(&b, simProfile)
	s.NoError(err)
	actual := b.String()
	expected := test.GoldenValueRaw(s.T(), path, []byte(actual))
	s.Equal(string(expected), actual)
}

func (s *RenderTestSuite) TestMarshal_Warrior() {
	// input, err := os.ReadFile("testdata/test_warrior_in.simc")
	// s.NoError(err)

	// res, err := NewProfile(bytes.NewReader(input))
	// s.NoError(err)
	// s.NotNil(res)

	// path, err := filepath.Abs("testdata/test_warrior_out.json")
	// s.NoError(err)
	// expected := test.GoldenValueJSON(s.T(), path, res)
	// s.Equal(expected, res)
}
