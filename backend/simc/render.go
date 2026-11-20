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

	res := []string{
		input.Character.render(),
		input.Base.render(),
	}

	for _, id := range slices.Sorted(maps.Keys(input.Profilesets)) {
		res = append(res, input.Profilesets[id].renderWithID(id))
	}

	res = append(res, input.Config.render())

	_, err := io.WriteString(w, strings.Join(res, "\n"))
	return err
}
