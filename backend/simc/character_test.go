package simc

import (
	"strings"
	"testing"

	"github.com/RobertsMJ/simc-cloud/backend/simc/internal"
	"github.com/stretchr/testify/suite"
)

type CharacterTestSuite struct {
	suite.Suite
}

func TestCharacterTestSuite(t *testing.T) {
	suite.Run(t, new(CharacterTestSuite))
}

func (suite *CharacterTestSuite) TestNewCharacter() {
	input := `
		demonhunter=Kestdh
		level=90
		race=night_elf
		region=us
		server=illidan
		role=attack
		professions=herbalism=88/mining=84
		spec=havoc
	`
	statements, err := internal.ParseStatements(strings.NewReader(input))
	suite.NoError(err)

	expected := Character{
		Class:  WowClassDH,
		Name:   "Kestdh",
		Level:  90,
		Role:   "attack",
		Spec:   "havoc",
		Server: "illidan",
		Region: "us",
		Race:   "night_elf",
	}

	actual, err := NewCharacter(statements)
	suite.NoError(err)
	suite.Equal(expected, actual)
}

func (suite *CharacterTestSuite) TestMarshalSimc() {
	c := Character{
		Class:  WowClassDH,
		Name:   `"Kestdh"`,
		Level:  90,
		Role:   "attack",
		Spec:   "havoc",
		Server: "illidan",
		Region: "us",
		Race:   "night_elf",
	}

	expected := []string{
		`demonhunter="Kestdh"`,
		"level=90",
		"race=night_elf",
		"region=us",
		"server=illidan",
		"role=attack",
		"spec=havoc",
	}

	actual, err := c.MarshalSimc()
	suite.NoError(err)
	suite.ElementsMatch(expected, strings.Split(actual, "\n"))
}
