package simc

import (
	"fmt"
	"io"
	"maps"
	"slices"
	"strings"
)

type SimInput struct {
	Character   Character          `json:"character"`
	Base        Loadout            `json:"base"`
	Profilesets map[string]Loadout `json:"profilesets"`
	Config      Config             `json:"config"`
}

func Render(w io.Writer, input SimInput) error {

	res := []string{}
	character, err := input.Character.MarshalSimc()
	if err != nil {
		return err
	}
	res = append(res, character)

	base, err := input.Base.MarshalSimc()
	if err != nil {
		return err
	}
	res = append(res, base)

	for _, id := range slices.Sorted(maps.Keys(input.Profilesets)) {
		profile := input.Profilesets[id]

		profileStr, err := profile.MarshalSimc()
		if err != nil {
			return err
		}

		profilesetLines := strings.Split(profileStr, "\n")
		for i, line := range profilesetLines {
			if i == 0 {
				profilesetLines[i] = fmt.Sprintf(`profileset."%s"=%s`, id, line)
			} else {
				profilesetLines[i] = fmt.Sprintf(`profileset."%s"+=%s`, id, line)
			}
		}

		res = append(res, strings.Join(profilesetLines, "\n"))
	}

	config, err := input.Config.MarshalSimc()
	if err != nil {
		return err
	}
	res = append(res, config)

	_, err = io.WriteString(w, strings.Join(res, "\n"))
	return err
}
