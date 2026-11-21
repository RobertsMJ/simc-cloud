package gearset_test

import (
	"encoding/json"
	"log/slog"
	"os"
	"testing"

	"github.com/RobertsMJ/simc-cloud/backend/gearset"
	"github.com/RobertsMJ/simc-cloud/backend/models"
	"github.com/RobertsMJ/simc-cloud/backend/simc"
	"github.com/stretchr/testify/suite"
)

type testCase struct {
	Base    simc.Loadout      `json:"base"`
	Options models.Candidates `json:"options"`
}

type GearsetTestSuite struct {
	suite.Suite

	underTest *gearset.Service
}

func TestGearsetTestSuite(t *testing.T) {
	suite.Run(t, new(GearsetTestSuite))
}

func (s *GearsetTestSuite) SetupTest() {
	s.underTest = gearset.NewService(slog.Default())
}

func (s *GearsetTestSuite) TestGenerateGearsets() {

	var testCase testCase
	// 2 trinkets, 1 talent option
	b, err := os.ReadFile("testdata/gearsets.json")
	s.NoError(err)
	err = json.Unmarshal(b, &testCase)
	s.NoError(err)

	gearsets := s.underTest.GenerateGearsets(testCase.Base, testCase.Options)
	s.Len(gearsets, 12)

}

func (s *GearsetTestSuite) TestGenerateGearsets_NoOptions() {

	var testCase testCase
	// 2 trinkets, 1 talent option
	b, err := os.ReadFile("testdata/gearsets.json")
	s.NoError(err)
	err = json.Unmarshal(b, &testCase)
	s.NoError(err)

	gearsets := s.underTest.GenerateGearsets(testCase.Base, models.Candidates{})
	s.Len(gearsets, 1)
}
