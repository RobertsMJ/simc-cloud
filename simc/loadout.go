package simc

import (
	"errors"

	"github.com/RobertsMJ/simc-cloud/simc/internal"
)

var ErrInvalidLoadoutStatement = errors.New("invalid loadout statement")

type Loadout struct {
	Talent        Talent
	Items         map[ItemSlot]Item
	OmniumTalents internal.IDValueList
}

type LoadoutOptions struct {
	Talents       []Talent
	Items         []Item
	OmniumTalents []internal.IDValueList
}

var (
	_ internal.ValueMarshaler   = Loadout{}
	_ internal.ValueUnmarshaler = (*Loadout)(nil)
	_ internal.ValueMarshaler   = LoadoutOptions{}
	_ internal.ValueUnmarshaler = (*LoadoutOptions)(nil)
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

func (l Loadout) MarshalSimcValue() (string, error) {
	panic("not implemented")
}

func (l Loadout) UnmarshalSimcValue(value string) error {
	panic("not implemented")
}

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
		var omniumTalents internal.IDValueList
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
		lo.Items = append(lo.Items, item)
		return nil
	}
	return ErrInvalidLoadoutStatement
}

func (lo LoadoutOptions) MarshalSimcValue() (string, error) {
	panic("not implemented")
}

func (lo LoadoutOptions) UnmarshalSimcValue(value string) error {
	panic("not implemented")
}
