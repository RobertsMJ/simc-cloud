package simc

import (
	"fmt"
	"strings"
)

var ErrInvalidTalentStatement = fmt.Errorf("invalid talent statement")

const (
	talentHeaderTemplate = "# Saved Loadout: %s"
	talentValueTemplate  = "talent=%s"
)

type Talent struct {
	Name  string
	Value string
}

var (
	_ StatementUnmarshaler = (*Talent)(nil)
	_ ValueMarshaler       = Talent{}
)

func (t *Talent) UnmarshalStatement(s statement) error {
	if s.Key != "talent" {
		return ErrInvalidTalentStatement
	}

	t.Name = strings.TrimPrefix(s.Comment, "Saved Loadout: ")
	t.Value = s.RawSimc
	return nil
}

func (t Talent) MarshalSimcValue() (string, error) {
	return strings.Join([]string{
		fmt.Sprintf(talentHeaderTemplate, t.Name),
		fmt.Sprintf(talentValueTemplate, t.Value),
	}, "\n"), nil
}
