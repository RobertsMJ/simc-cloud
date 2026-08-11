package simc

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type ItemTestSuite struct {
	suite.Suite
}

func TestItemTestSuite(t *testing.T) {
	suite.Run(t, new(ItemTestSuite))
}

func (s *ItemTestSuite) TestUnmarshalStatement() {
	in := `head=,id=251109,enchant_id=8017,bonus_id=13440/6652/12667/13577/12699/12806,foo=bar`
	stmt := newStatement(in, 0)

	var item Item
	err := UnmarshalStatement(stmt, &item)
	s.NoError(err)
}
