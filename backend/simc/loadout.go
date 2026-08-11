package simc

import "errors"

var ErrInvalidLoadoutStatement = errors.New("invalid loadout statement")

type Loadout struct {
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
	_ ValueMarshaler   = Loadout{}
	_ ValueUnmarshaler = (*Loadout)(nil)
	_ ValueMarshaler   = LoadoutOptions{}
	_ ValueUnmarshaler = (*LoadoutOptions)(nil)
)

func (l Loadout) UnmarshalStatement(s statement) error {
	if s.Key == "talent" {
		return l.Talent.UnmarshalStatement(s)
	}
	if s.Key == "omnium_talents" {
		return l.OmniumTalents.UnmarshalStatement(s)
	}
	if itemSlot, ok := ItemSlotFromString(s.Key); ok {
		var item Item
		if err := UnmarshalStatement(s, &item); err != nil {
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

func (lo *LoadoutOptions) UnmarshalStatement(s statement) error {
	if s.Key == "talent" {
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
		if err := UnmarshalStatement(s, &item); err != nil {
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
