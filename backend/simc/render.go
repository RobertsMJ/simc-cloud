package simc

import (
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

// Render writes a valid simulation input to the provided writer.
func Render(w io.Writer, input SimInput) error {

	if _, err := io.WriteString(w, strings.Join([]string{
		input.Character.render(),
		input.Base.render(),
	}, "\n")+"\n"); err != nil {
		return err
	}

	for _, id := range slices.Sorted(maps.Keys(input.Profilesets)) {
		if _, err := io.WriteString(w, input.Profilesets[id].renderWithID(id)+"\n"); err != nil {
			return err
		}
	}

	_, err := io.WriteString(w, input.Config.render())
	return err
}
