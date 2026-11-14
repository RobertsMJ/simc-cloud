package simc

import (
	"errors"
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/RobertsMJ/simc-cloud/backend/simc/internal"
)

var ErrInvalidLoadoutStatement = errors.New("invalid loadout statement")

type Loadout struct {
	ID    string            `json:"id,omitempty"`
	Items map[ItemSlot]Item `json:"items"`

	Talent        Talent      `json:"talents"`
	OmniumTalents IDValueList `json:"omnium_talents,omitempty" simc:"omnium_talents"`
}

var (
	_ internal.ValueMarshaler = Loadout{}
)

func (l *Loadout) UnmarshalStatement(s internal.Statement) error {
	if s.Key == "talents" {
		return l.Talent.UnmarshalStatement(s)
	}
	if s.Key == "omnium_talents" {
		return l.OmniumTalents.UnmarshalStatement(s)
	}
	if itemSlot, ok := ItemSlotFromString(s.Key); ok {
		var item Item
		if err := item.UnmarshalStatement(s); err != nil {
			return err
		}
		if l.Items == nil {
			l.Items = make(map[ItemSlot]Item)
		}
		l.Items[itemSlot] = item
		return nil
	}
	return ErrInvalidLoadoutStatement
}

func (l Loadout) MarshalSimc() (string, error) {
	res := []string{}

	for _, slot := range slices.Sorted(maps.Keys(l.Items)) {
		item := l.Items[slot]
		str, err := item.MarshalSimc()
		if err != nil {
			return "", err
		}
		res = append(res, str)
	}

	fields, err := internal.MarshalFields(l)
	if err != nil {
		return "", err
	}
	res = append(res, fields...)
	if l.ID != "" {
		prefix := l.ID
		for i, str := range res {
			if i == 0 {
				res[i] = fmt.Sprintf(`profileset."%s"=%s`, prefix, str)
			} else {
				res[i] = fmt.Sprintf(`profileset."%s"+=%s`, prefix, str)
			}
		}
	}

	return strings.Join(res, "\n"), nil
}

type LoadoutOptions struct {
	Talents       []Talent             `json:"talents,omitempty"`
	Items         map[SlotGroup][]Item `json:"items,omitempty"`
	OmniumTalents []IDValueList        `json:"omnium_talents,omitempty"`
}

var (
	_ internal.ValueMarshaler = LoadoutOptions{}
)

func (lo *LoadoutOptions) UnmarshalStatement(s internal.Statement) error {
	if s.Key == "talents" {
		var talent Talent
		if err := talent.UnmarshalStatement(s); err != nil {
			return err
		}
		lo.Talents = append(lo.Talents, talent)
		return nil
	}
	if s.Key == "omnium_talents" {
		var omniumTalents IDValueList
		if err := omniumTalents.UnmarshalStatement(s); err != nil {
			return err
		}
		lo.OmniumTalents = append(lo.OmniumTalents, omniumTalents)
		return nil
	}
	if itemSlot, ok := ItemSlotFromString(s.Key); ok {
		var item Item
		if err := item.UnmarshalStatement(s); err != nil {
			return err
		}
		item.Slot = itemSlot
		slotGroup := itemSlot.GetSlotGroup()
		lo.Items[slotGroup] = append(lo.Items[slotGroup], item)
		return nil
	}
	return ErrInvalidLoadoutStatement
}

// TODO
func (lo LoadoutOptions) MarshalSimc() (string, error) {
	return "", nil
}
