package simc_test

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/RobertsMJ/simc-cloud-backend/simc"
)

type SimCInputTestSuite struct {
	suite.Suite
}

func TestSimCInputTestSuite(t *testing.T) {
	suite.Run(t, new(SimCInputTestSuite))
}

// Read in a simc input file and test the Unmarshal function to ensure it correctly parses the character, equipment, and options sections into the SimCInput struct
func (s *SimCInputTestSuite) TestUnmarshalSimCInput() {
	input, err := os.ReadFile("./test-data/sample.simc")
	s.NoError(err)

	var simcInput simc.Input
	err = simcInput.UnmarshalSimC(input)
	s.NoError(err)

	s.Equal("Kestdh", simcInput.Character.CharName)
	s.NotNil(simcInput.Equipment)
	s.NotNil(simcInput.Equipment[simc.EquipmentSlotHead])
	head := simcInput.Equipment[simc.EquipmentSlotHead]
	s.Equal(251109, head.ID)
}

func (s *SimCInputTestSuite) TestMarshalSimCInput() {
	input, err := os.ReadFile("./test-data/sample.simc")
	s.NoError(err)
	// Expected is the same input string with blank lines and comments removed, since those are not preserved in the struct representation
	lines := strings.Split(string(input), "\n")
	var expectedLines []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") {
			expectedLines = append(expectedLines, line)
		}
	}

	var simcInput simc.Input
	err = simcInput.UnmarshalSimC(input)
	s.NoError(err)

	marshaled, err := simcInput.MarshalSimC()
	s.NoError(err)

	actualLines := strings.Split(string(marshaled), "\n")

	s.ElementsMatch(expectedLines, actualLines)
}
