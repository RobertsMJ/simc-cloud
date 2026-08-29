package internal_test

import (
	"testing"

	"github.com/RobertsMJ/simc-cloud/simc/internal"
	"github.com/stretchr/testify/suite"
)

type IDListTestSuite struct {
	suite.Suite
}

func TestIDListTestSuite(t *testing.T) {
	suite.Run(t, new(IDListTestSuite))
}

func (s *IDListTestSuite) TestMarshalSimcValue() {
	// Test implementation
	idList := internal.IDList{1, 2, 3}
	result, err := idList.MarshalSimcValue()
	s.NoError(err)
	s.Equal("1/2/3", result)
}

func (s *IDListTestSuite) TestUnmarshalSimcValue() {
	// Test implementation
	idList := internal.IDList{}
	err := idList.UnmarshalSimcValue("1/2/3")
	s.NoError(err)
	s.Equal(internal.IDList{1, 2, 3}, idList)
}

func (s *IDListTestSuite) TestMarshalSimcValue_Empty() {
	// Test implementation
	idList := internal.IDList{}
	result, err := idList.MarshalSimcValue()
	s.NoError(err)
	s.Equal("", result)
}

func (s *IDListTestSuite) TestUnmarshalSimcValue_Empty() {
	// Test implementation
	idList := internal.IDList{}
	err := idList.UnmarshalSimcValue("")
	s.NoError(err)
	s.Equal(internal.IDList{}, idList)
}
