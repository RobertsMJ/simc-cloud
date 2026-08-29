package simc

import (
	"fmt"
	"strings"

	"github.com/RobertsMJ/simc-cloud/simc/internal"
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
	_ internal.StatementUnmarshaler = (*Talent)(nil)
	_ internal.ValueMarshaler       = Talent{}
)

func (t *Talent) UnmarshalStatement(s internal.Statement) error {
	if s.Key != "talents" {
		return ErrInvalidTalentStatement
	}

	t.Name = strings.TrimPrefix(s.Comment, "Saved Loadout: ")
	t.Value = s.Value
	return nil
}

func (t Talent) MarshalSimcValue() (string, error) {
	return strings.Join([]string{
		fmt.Sprintf(talentHeaderTemplate, t.Name),
		fmt.Sprintf(talentValueTemplate, t.Value),
	}, "\n"), nil
}
