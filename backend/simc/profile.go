package simc

import (
	"strings"

	"github.com/RobertsMJ/simc-cloud/backend/simc/internal"
)

type Profile struct {
	Character Character `json:"character"`
	Loadout   Loadout   `json:"loadout"`
}

var (
	_ internal.ValueMarshaler = Profile{}
)

func (p Profile) MarshalSimc() (string, error) {
	character, err := p.Character.MarshalSimc()
	if err != nil {
		return "", err
	}
	loadout, err := p.Loadout.MarshalSimc()
	if err != nil {
		return "", err
	}
	return strings.Join([]string{character, loadout}, "\n"), nil
}
