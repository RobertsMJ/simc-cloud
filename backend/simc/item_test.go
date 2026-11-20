package simc

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ItemTestSuite struct {
	suite.Suite
}

func TestItemTestSuite(t *testing.T) {
	suite.Run(t, new(ItemTestSuite))
}

func (s *ItemTestSuite) TestNewItem() {
	in := `
	# New Helmet
	head=,id=251109,enchant_id=8017,bonus_id=13440/6652/12667/13577/12699/12806,foo=bar
	`
	stmt, _ := parseStatements(strings.NewReader(in))
	s.Len(stmt, 1)

	expected := Item{
		Slot:      SlotGroupHead,
		Name:      "New Helmet",
		ID:        251109,
		EnchantID: IDList{8017},
		BonusID:   IDList{13440, 6652, 12667, 13577, 12699, 12806},
	}

	var actual Item
	actual, err := newItem(stmt[0])
	s.NoError(err)
	s.Equal(expected, actual)
}

func (s *ItemTestSuite) TestRender() {
	in := Item{
		Slot:      SlotGroupHead,
		Name:      "New Helmet",
		ID:        251109,
		EnchantID: IDList{8017},
		BonusID:   IDList{13440, 6652, 12667, 13577, 12699, 12806},
	}
	expected := `head=,id=251109,enchant_id=8017,bonus_id=13440/6652/12667/13577/12699/12806`
	actualStr := in.render(ItemSlotHead)
	s.Equal(expected, actualStr)
}
