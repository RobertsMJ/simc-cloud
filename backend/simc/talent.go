package simc

import (
	"fmt"
	"strings"

	"github.com/RobertsMJ/simc-cloud/backend/simc/internal"
)

var (
	ErrInvalidTalentStatement       = fmt.Errorf("invalid talent statement")
	ErrInvalidOmniumTalentStatement = fmt.Errorf("invalid omnium talent statement")
)

const (
	talentHeaderTemplate = "# Saved Loadout: %s"
	talentValueTemplate  = "talents=%s"
	talentKey            = "talents"
	talentCommentPrefix  = "Saved Loadout: "
	omniumTalentKey      = "omnium_talents"
	omniumTalentTemplate = "omnium_talents=%s"
)

type Talent struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func newTalent(s internal.Statement) (Talent, error) {
	if s.Key != talentKey {
		return Talent{}, ErrInvalidTalentStatement
	}

	return Talent{
		Name:  strings.TrimPrefix(s.Comment, talentCommentPrefix),
		Value: s.Value,
	}, nil
}

func (t Talent) MarshalSimc() (string, error) {
	return fmt.Sprintf(talentValueTemplate, t.Value), nil
}

type OmniumTalents = IDValueList

func newOmniumTalents(s internal.Statement) (OmniumTalents, error) {
	if s.Key != omniumTalentKey {
		return nil, ErrInvalidOmniumTalentStatement
	}

	return NewIDValueList(s.Value)
}
