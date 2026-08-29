package internal_test

import (
	"testing"

	"github.com/RobertsMJ/simc-cloud/simc/internal"
	"github.com/stretchr/testify/suite"
)

type IDValueTestSuite struct {
	suite.Suite
}

func TestIDValueTestSuite(t *testing.T) {
	suite.Run(t, new(IDValueTestSuite))
}

func (s *IDValueTestSuite) TestMarshalSimcValue_IDValue() {
	// Test implementation
	idValue := internal.IDValue{ID: 1, Value: 2}
	result, err := idValue.MarshalSimcValue()
	s.NoError(err)
	s.Equal("1:2", result)
}

func (s *IDValueTestSuite) TestUnmarshalSimcValue_IDValue() {
	// Test implementation
	idValue := internal.IDValue{}
	err := idValue.UnmarshalSimcValue("1:2")
	s.NoError(err)
	s.Equal(internal.IDValue{ID: 1, Value: 2}, idValue)
}

func (s *IDValueTestSuite) TestMarshalSimcValue_IDValueList() {
	// Test implementation
	idList := internal.IDValueList{{ID: 1, Value: 2}, {ID: 3, Value: 4}}
	result, err := idList.MarshalSimcValue()
	s.NoError(err)
	s.Equal("1:2/3:4", result)
}

func (s *IDValueTestSuite) TestUnmarshalSimcValue_IDValueList() {
	// Test implementation
	idList := internal.IDValueList{}
	err := idList.UnmarshalSimcValue("1:2/3:4")
	s.NoError(err)
	s.Equal(internal.IDValueList{{ID: 1, Value: 2}, {ID: 3, Value: 4}}, idList)
}
