package simc

import (
	"fmt"
	"strings"
)

var (
	ErrInvalidTalentStatement       = fmt.Errorf("invalid talent statement")
	ErrInvalidOmniumTalentStatement = fmt.Errorf("invalid omnium talent statement")
)

const (
	talentKey           = "talents"
	talentCommentPrefix = "Saved Loadout: "
	omniumTalentKey     = "omnium_talents"
)

type Talent struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func newTalent(name string, value string) (Talent, error) {
	return Talent{
		Name:  strings.TrimPrefix(name, talentCommentPrefix),
		Value: value,
	}, nil
}

func (t Talent) render() string {
	return fmt.Sprintf("talents=%s", t.Value)
}

type OmniumTalents = IDValueList

func newOmniumTalents(value string) (OmniumTalents, error) {
	return newIDValueList(value)
}
