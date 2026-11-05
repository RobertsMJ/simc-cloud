package simc

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type IDValueTestSuite struct {
	suite.Suite
}

func TestIDValueTestSuite(t *testing.T) {
	suite.Run(t, new(IDValueTestSuite))
}

func (s *IDValueTestSuite) TestMarshalSimc_IDValue() {
	// Test implementation
	idValue := IDValue{ID: 1, Value: 2}
	result, err := idValue.MarshalSimc()
	s.NoError(err)
	s.Equal("1:2", result)
}

func (s *IDValueTestSuite) TestUnmarshalSimc_IDValue() {
	// Test implementation
	idValue := IDValue{}
	err := idValue.UnmarshalSimc("1:2")
	s.NoError(err)
	s.Equal(IDValue{ID: 1, Value: 2}, idValue)
}

func (s *IDValueTestSuite) TestMarshalSimc_IDValueList() {
	// Test implementation
	idList := IDValueList{{ID: 1, Value: 2}, {ID: 3, Value: 4}}
	result, err := idList.MarshalSimc()
	s.NoError(err)
	s.Equal("1:2/3:4", result)
}

func (s *IDValueTestSuite) TestUnmarshalSimc_IDValueList() {
	// Test implementation
	idList := IDValueList{}
	err := idList.UnmarshalSimc("1:2/3:4")
	s.NoError(err)
	s.Equal(IDValueList{{ID: 1, Value: 2}, {ID: 3, Value: 4}}, idList)
}
