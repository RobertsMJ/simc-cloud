package simc_test

import (
	"strings"
	"testing"

	"github.com/RobertsMJ/simc-cloud/simc"
	"github.com/RobertsMJ/simc-cloud/simc/internal"
	"github.com/stretchr/testify/suite"
)

type ItemTestSuite struct {
	suite.Suite
}

func TestItemTestSuite(t *testing.T) {
	suite.Run(t, new(ItemTestSuite))
}

func (s *ItemTestSuite) TestUnmarshalStatement() {
	in := `
	# New Helmet
	head=,id=251109,enchant_id=8017,bonus_id=13440/6652/12667/13577/12699/12806,foo=bar
	`
	stmt, _ := internal.ParseStatements(strings.NewReader(in))
	s.Len(stmt, 1)

	expected := simc.Item{
		Slot:      simc.ItemSlotHead,
		SlotGroup: "",
		Name:      "New Helmet",
		ID:        251109,
		EnchantID: internal.IDList{8017},
		BonusID:   internal.IDList{13440, 6652, 12667, 13577, 12699, 12806},
	}

	var actual simc.Item
	err := actual.UnmarshalStatement(stmt[0])
	s.NoError(err)
	s.Equal(expected, actual)
}
