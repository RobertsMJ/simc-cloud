package internal_test

import (
	"errors"
	"testing"

	"github.com/RobertsMJ/simc-cloud/backend/simc/internal"
	"github.com/stretchr/testify/suite"
)

var errMarshal = errors.New("marshal failed")

type upper string

func (u upper) MarshalSimc() (string, error) { return "UP:" + string(u), nil }

type failing struct{}

func (failing) MarshalSimc() (string, error) { return "", errMarshal }

type MarshalFieldsTestSuite struct {
	suite.Suite
}

func TestMarshalFieldsTestSuite(t *testing.T) {
	suite.Run(t, new(MarshalFieldsTestSuite))
}

func (s *MarshalFieldsTestSuite) TestNilPointerInput() {
	var in *struct{}
	_, err := internal.MarshalFields(in)
	s.ErrorIs(err, internal.ErrNilMarshal)
}

func (s *MarshalFieldsTestSuite) TestValueAndPointerInputMatch() {
	in := struct {
		ID int `simc:"id"`
	}{ID: 7}

	fromValue, err := internal.MarshalFields(in)
	s.Require().NoError(err)
	fromPointer, err := internal.MarshalFields(&in)
	s.Require().NoError(err)

	s.Equal([]string{"id=7"}, fromValue)
	s.Equal(fromValue, fromPointer)
}

func (s *MarshalFieldsTestSuite) TestPlainValuesInFieldOrder() {
	in := struct {
		Name  string  `simc:"name"`
		Level int     `simc:"level"`
		Ratio float64 `simc:"ratio"`
	}{Name: "kest", Level: 90, Ratio: 0.5}

	res, err := internal.MarshalFields(in)
	s.Require().NoError(err)
	s.Equal([]string{"name=kest", "level=90", "ratio=0.5"}, res)
}

func (s *MarshalFieldsTestSuite) TestSkipsUntaggedAndIgnoredFields() {
	in := struct {
		Untagged string
		Ignored  string `simc:"-"`
		ID       int    `simc:"id"`
	}{Untagged: "a", Ignored: "b", ID: 1}

	res, err := internal.MarshalFields(in)
	s.Require().NoError(err)
	s.Equal([]string{"id=1"}, res)
}

func (s *MarshalFieldsTestSuite) TestSkipsUnexportedTaggedField() {
	in := struct {
		hidden int `simc:"hidden"`
		ID     int `simc:"id"`
	}{hidden: 1, ID: 2}

	res, err := internal.MarshalFields(in)
	s.Require().NoError(err)
	s.Equal([]string{"id=2"}, res)
}

func (s *MarshalFieldsTestSuite) TestZeroValueWithoutOmitempty() {
	in := struct {
		ID int `simc:"id"`
	}{}

	res, err := internal.MarshalFields(in)
	s.Require().NoError(err)
	s.Equal([]string{"id=0"}, res)
}

func (s *MarshalFieldsTestSuite) TestOmitemptySkipsZeroValue() {
	in := struct {
		ID int `simc:"id,omitempty"`
	}{}

	res, err := internal.MarshalFields(in)
	s.Require().NoError(err)
	s.Empty(res)
}

func (s *MarshalFieldsTestSuite) TestNilPointerWithoutOmitempty() {
	in := struct {
		Tuning *int `simc:"content_tuning"`
	}{}

	res, err := internal.MarshalFields(in)
	s.Require().NoError(err)
	s.Empty(res)
}

func (s *MarshalFieldsTestSuite) TestNilPointerWithOmitempty() {
	in := struct {
		Tuning *int `simc:"content_tuning,omitempty"`
	}{}

	res, err := internal.MarshalFields(in)
	s.Require().NoError(err)
	s.Empty(res)
}

func (s *MarshalFieldsTestSuite) TestPointerIsDereferenced() {
	tuning := 3615
	in := struct {
		Tuning *int `simc:"content_tuning"`
	}{Tuning: &tuning}

	res, err := internal.MarshalFields(in)
	s.Require().NoError(err)
	s.Equal([]string{"content_tuning=3615"}, res)
}

func (s *MarshalFieldsTestSuite) TestOmitemptyKeepsPointerToZero() {
	zero := 0
	in := struct {
		Tuning *int `simc:"content_tuning,omitempty"`
	}{Tuning: &zero}

	res, err := internal.MarshalFields(in)
	s.Require().NoError(err)
	s.Equal([]string{"content_tuning=0"}, res)
}

func (s *MarshalFieldsTestSuite) TestOmitemptySkipsNilSlice() {
	in := struct {
		Gems []int `simc:"gem_id,omitempty"`
	}{}

	res, err := internal.MarshalFields(in)
	s.Require().NoError(err)
	s.Empty(res)
}

func (s *MarshalFieldsTestSuite) TestOmitemptySkipsEmptySlice() {
	in := struct {
		Gems []int `simc:"gem_id,omitempty"`
	}{Gems: []int{}}

	res, err := internal.MarshalFields(in)
	s.Require().NoError(err)
	s.Empty(res)
}

func (s *MarshalFieldsTestSuite) TestUsesValueMarshaler() {
	in := struct {
		Name upper `simc:"name"`
	}{Name: "kest"}

	res, err := internal.MarshalFields(in)
	s.Require().NoError(err)
	s.Equal([]string{"name=UP:kest"}, res)
}

func (s *MarshalFieldsTestSuite) TestUsesValueMarshalerBehindPointer() {
	name := upper("kest")
	in := struct {
		Name *upper `simc:"name"`
	}{Name: &name}

	res, err := internal.MarshalFields(in)
	s.Require().NoError(err)
	s.Equal([]string{"name=UP:kest"}, res)
}

func (s *MarshalFieldsTestSuite) TestNilMarshalerPointerIsSkipped() {
	in := struct {
		Name *upper `simc:"name"`
	}{}

	res, err := internal.MarshalFields(in)
	s.Require().NoError(err)
	s.Empty(res)
}

func (s *MarshalFieldsTestSuite) TestMarshalerErrorIncludesKey() {
	in := struct {
		Broken failing `simc:"broken"`
	}{}

	_, err := internal.MarshalFields(in)
	s.ErrorIs(err, errMarshal)
	s.ErrorContains(err, "broken")
}
