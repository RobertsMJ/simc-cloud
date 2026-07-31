package simc

import "fmt"

type Profile struct {
	Character Character `json:"character"`
	Loadout   Loadout   `json:"loadout"`
	Config    SimConfig `json:"config"`
}

var (
	_ ValueMarshaler   = Profile{}
	_ ValueUnmarshaler = (*Profile)(nil)
)

func (p Profile) MarshalSimcValue() (string, error) {
	character, err := p.Character.MarshalSimcValue()
	if err != nil {
		return "", err
	}
	loadout, err := p.Loadout.MarshalSimcValue()
	if err != nil {
		return "", err
	}
	cfg, err := p.Config.MarshalSimcValue()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(`
%s
%s
%s
	`, character, loadout, cfg), nil
}

func (p Profile) UnmarshalSimcValue(value string) error {
	panic("not implemented")
}
