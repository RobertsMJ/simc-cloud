package simc

type Profile struct {
	Character Character
	Loadout   Loadout
	Config    SimConfig
}

var (
	_ ValueMarshaler   = Profile{}
	_ ValueUnmarshaler = (*Profile)(nil)
)

func (p Profile) MarshalSimcValue() (string, error) {
	panic("not implemented")
}

func (p Profile) UnmarshalSimcValue(value string) error {
	panic("not implemented")
}
