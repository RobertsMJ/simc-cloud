package simc_test

import (
	"strings"
	"testing"

	"github.com/RobertsMJ/simc-cloud/simc"
	"github.com/RobertsMJ/simc-cloud/simc/internal"
	"github.com/stretchr/testify/suite"
)

type CharacterTestSuite struct {
	suite.Suite
}

func TestCharacterTestSuite(t *testing.T) {
	suite.Run(t, new(CharacterTestSuite))
}

func (suite *CharacterTestSuite) TestUnmarshalStatements() {
	input := `
		demonhunter="Kestdh"
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

	expected := simc.Character{
		Class:  simc.WowClassDH,
		Name:   `"Kestdh"`,
		Level:  90,
		Role:   "attack",
		Spec:   "havoc",
		Server: "illidan",
		Region: "us",
		Race:   "night_elf",
	}

	var actual simc.Character
	for _, stmt := range statements {
		err := actual.UnmarshalStatement(stmt)
		suite.NoError(err)
	}

	suite.Equal(expected, actual)
}
