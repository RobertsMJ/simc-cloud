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
	result, err := idList.MarshalSimc()
	s.NoError(err)
	s.Equal("1/2/3", result)
}

func (s *IDListTestSuite) TestUnmarshalSimc() {
	// Test implementation
	idList := IDList{}
	err := idList.UnmarshalSimc("1/2/3")
	s.NoError(err)
	s.Equal(IDList{1, 2, 3}, idList)
}

func (s *IDListTestSuite) TestMarshalSimc_Empty() {
	// Test implementation
	idList := IDList{}
	result, err := idList.MarshalSimc()
	s.NoError(err)
	s.Equal("", result)
}

func (s *IDListTestSuite) TestUnmarshalSimc_Empty() {
	// Test implementation
	idList := IDList{}
	err := idList.UnmarshalSimc("")
	s.NoError(err)
	s.Equal(IDList{}, idList)
}
