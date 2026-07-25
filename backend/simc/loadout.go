package simc

type Loadout struct {
	Talent        Talent
	Items         map[ItemSlot]Item
	OmniumTalents IDValueList
}

type LoadoutOptions struct {
	Talents       []Talent
	Items         map[SlotGroup][]Item
	OmniumTalents IDValueList
}

var (
	_ ValueMarshaler   = Loadout{}
	_ ValueUnmarshaler = (*Loadout)(nil)
	_ ValueMarshaler   = LoadoutOptions{}
	_ ValueUnmarshaler = (*LoadoutOptions)(nil)
)

func (l Loadout) MarshalSimcValue() (string, error) {
	panic("not implemented")
}

func (l Loadout) UnmarshalSimcValue(value string) error {
	panic("not implemented")
}

func (lo LoadoutOptions) MarshalSimcValue() (string, error) {
	panic("not implemented")
}

func (lo LoadoutOptions) UnmarshalSimcValue(value string) error {
	panic("not implemented")
}
