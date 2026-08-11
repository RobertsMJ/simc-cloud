package simc

import "strings"

type Profile struct {
	Character Character `json:"character"`
	Loadout   Loadout   `json:"loadout"`
}

var (
	_ ValueMarshaler = Profile{}
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
	return strings.Join([]string{character, loadout}, "\n"), nil
}
