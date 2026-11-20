package simc

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type IDListTestSuite struct {
	suite.Suite
}

func TestIDListTestSuite(t *testing.T) {
	suite.Run(t, new(IDListTestSuite))
}

func (s *IDListTestSuite) TestMarshalSimc() {
	// Test implementation
	idList := IDList{1, 2, 3}
	result := idList.render()
	s.Equal("1/2/3", result)
}

func (s *IDListTestSuite) TestNewIDList() {
	// Test implementation
	idList, err := newIDList("1/2/3")
	s.NoError(err)
	s.Equal(IDList{1, 2, 3}, idList)
}

func (s *IDListTestSuite) TestMarshalSimc_Empty() {
	// Test implementation
	idList := IDList{}
	result := idList.render()
	s.Equal("", result)
}

func (s *IDListTestSuite) TestUnmarshalSimc_Empty() {
	// Test implementation
	idList, err := newIDList("")
	s.NoError(err)
	s.Equal(IDList{}, idList)
}
