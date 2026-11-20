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

func newTalent(s statement) (Talent, error) {
	if s.Key != talentKey {
		return Talent{}, ErrInvalidTalentStatement
	}

	return Talent{
		Name:  strings.TrimPrefix(s.Comment, talentCommentPrefix),
		Value: s.Value,
	}, nil
}

func (t Talent) render() string {
	return fmt.Sprintf("talents=%s", t.Value)
}

type OmniumTalents struct {
	Value IDValueList `json:"value"`
}

func newOmniumTalents(s statement) (OmniumTalents, error) {
	if s.Key != omniumTalentKey {
		return OmniumTalents{}, ErrInvalidOmniumTalentStatement
	}

	v, err := newIDValueList(s.Value)
	if err != nil {
		return OmniumTalents{}, err
	}

	return OmniumTalents{
		Value: v,
	}, nil
}

func (o OmniumTalents) render() string {
	return fmt.Sprintf("omnium_talents=%s", o.Value.render())
}
