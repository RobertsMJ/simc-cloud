package simc

import (
	"errors"

	"github.com/RobertsMJ/simc-cloud/backend/simc/internal"
)

var ErrInvalidLoadoutStatement = errors.New("invalid loadout statement")

type Loadout struct {
	ID            string
	Talent        Talent
	Items         map[ItemSlot]Item
	OmniumTalents IDValueList
}

type LoadoutOptions struct {
	Talents       []Talent
	Items         []Item
	OmniumTalents []IDValueList
}

var (
	_ internal.ValueMarshaler   = Loadout{}
	_ internal.ValueUnmarshaler = (*Loadout)(nil)
	_ internal.ValueMarshaler   = LoadoutOptions{}
	_ internal.ValueUnmarshaler = (*LoadoutOptions)(nil)
)

func (l *Loadout) unmarshalStatement(s internal.Statement) error {
	if s.Key == "talents" {
		return l.Talent.unmarshalStatement(s)
	}
	if s.Key == "omnium_talents" {
		return l.OmniumTalents.unmarshalStatement(s)
	}
	if itemSlot, ok := ItemSlotFromString(s.Key); ok {
		var item Item
		if err := item.unmarshalStatement(s); err != nil {
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
	panic("not implemented")
}

func (l Loadout) UnmarshalSimc(value string) error {
	panic("not implemented")
}

func (lo *LoadoutOptions) unmarshalStatement(s internal.Statement) error {
	if s.Key == "talents" {
		var talent Talent
		if err := talent.unmarshalStatement(s); err != nil {
			return err
		}
		lo.Talents = append(lo.Talents, talent)
		return nil
	}
	if s.Key == "omnium_talents" {
		var omniumTalents IDValueList
		if err := omniumTalents.unmarshalStatement(s); err != nil {
			return err
		}
		lo.OmniumTalents = append(lo.OmniumTalents, omniumTalents)
		return nil
	}
	if itemSlot, ok := ItemSlotFromString(s.Key); ok {
		var item Item
		if err := item.unmarshalStatement(s); err != nil {
			return err
		}
		item.Slot = itemSlot
		lo.Items = append(lo.Items, item)
		return nil
	}
	return ErrInvalidLoadoutStatement
}

func (lo LoadoutOptions) MarshalSimc() (string, error) {
	panic("not implemented")
}

func (lo LoadoutOptions) UnmarshalSimc(value string) error {
	panic("not implemented")
}
