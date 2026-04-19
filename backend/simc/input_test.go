package simc_test

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/RobertsMJ/simc-cloud-backend/simc"
)

// Read in a simc input file and test the Unmarshal function to ensure it correctly parses the character, equipment, and options sections into the SimCInput struct
func TestUnmarshalSimCInput(t *testing.T) {
	input, err := os.ReadFile("./test-data/sample.simc")
	if err != nil {
		t.Fatal(err)
	}

	simcInput, err := simc.Unmarshal(input)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "Kestdh", simcInput.Character.CharName)
	assert.NotNil(t, simcInput.Equipment)
	assert.NotNil(t, simcInput.Equipment[simc.EquipmentSlotHead])
	head := simcInput.Equipment[simc.EquipmentSlotHead]
	assert.Equal(t, 251109, head.ID)
}

func TestMarshalSimCInput(t *testing.T) {
	input, err := os.ReadFile("./test-data/sample.simc")
	if err != nil {
		t.Fatal(err)
	}
	// Expected is the same input string with blank lines and comments removed, since those are not preserved in the struct representation
	lines := strings.Split(string(input), "\n")
	var expectedLines []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") {
			expectedLines = append(expectedLines, line)
		}
	}

	simcInput, err := simc.Unmarshal(input)
	if err != nil {
		t.Fatal(err)
	}

	marshaled, err := simcInput.MarshalSimC()
	if err != nil {
		t.Fatal(err)
	}

	actualLines := strings.Split(string(marshaled), "\n")

	assert.ElementsMatch(t, expectedLines, actualLines)
}
