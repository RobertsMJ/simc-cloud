package simc_test

import (
	"testing"

	"github.com/RobertsMJ/simc-cloud-backend/simc"
	"github.com/stretchr/testify/suite"
)

type ProfessionTest struct {
	suite.Suite
}

func TestProfession(t *testing.T) {
	suite.Run(t, new(ProfessionTest))
}

func (s *ProfessionTest) TestMarshalSimcValue_Profession() {
	prof := simc.Profession{
		ID:    "herbalism",
		Value: 88,
	}
	expected := "herbalism=88"
	result, err := prof.MarshalSimcValue()
	s.NoError(err)
	s.Equal(expected, result)
}

func (s *ProfessionTest) TestUnmarshalSimcValue_Profession() {
	var prof simc.Profession
	err := prof.UnmarshalSimcValue("herbalism=88")
	s.NoError(err)
	s.Equal("herbalism", prof.ID)
	s.Equal(88, prof.Value)
}

func (s *ProfessionTest) TestMarshalSimcValue_Professions() {
	professions := simc.Professions{
		{ID: "herbalism", Value: 88},
		{ID: "alchemy", Value: 90},
	}
	expected := "herbalism=88/alchemy=90"
	result, err := professions.MarshalSimcValue()
	s.NoError(err)
	s.Equal(expected, result)
}

func (s *ProfessionTest) TestUnmarshalSimcValue_Professions() {
	var professions simc.Professions
	err := professions.UnmarshalSimcValue("herbalism=88/alchemy=90")
	s.NoError(err)
	s.Equal(2, len(professions))
	s.Equal("herbalism", professions[0].ID)
	s.Equal(88, professions[0].Value)
	s.Equal("alchemy", professions[1].ID)
	s.Equal(90, professions[1].Value)
}
