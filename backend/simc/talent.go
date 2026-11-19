package simc

import (
	"fmt"
	"strings"

	"github.com/RobertsMJ/simc-cloud/backend/simc/internal"
)

var ErrInvalidTalentStatement = fmt.Errorf("invalid talent statement")

const (
	talentHeaderTemplate = "# Saved Loadout: %s"
	talentValueTemplate  = "talents=%s"
	talentKey            = "talents"
	talentCommentPrefix  = "Saved Loadout: "
)

type Talent struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

var (
	_ internal.ValueMarshaler = Talent{}
)

type OmniumTalents = IDValueList

func (t *Talent) UnmarshalStatement(s internal.Statement) error {
	if s.Key != talentKey {
		return ErrInvalidTalentStatement
	}

	t.Name = strings.TrimPrefix(s.Comment, talentCommentPrefix)
	t.Value = s.Value
	return nil
}

func (t Talent) MarshalSimc() (string, error) {
	return strings.Join([]string{
		fmt.Sprintf(talentHeaderTemplate, t.Name),
		fmt.Sprintf(talentValueTemplate, t.Value),
	}, "\n"), nil
}
