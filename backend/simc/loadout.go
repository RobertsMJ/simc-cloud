package simc

import (
	"errors"
	"maps"
	"slices"
	"strings"
)

var ErrInvalidLoadoutStatement = errors.New("invalid loadout statement")

// Loadout represents a character's active gear and talent configuration.
type Loadout struct {
	Items         map[ItemSlot]Item `json:"items"`
	Talent        Talent            `json:"talents"`
	OmniumTalents OmniumTalents     `json:"omnium_talents,omitempty"`
}

func (l Loadout) MarshalSimc() (string, error) {
	res := []string{}

	for _, slot := range slices.Sorted(maps.Keys(l.Items)) {
		item := l.Items[slot]
		str, err := item.MarshalSimc()
		if err != nil {
			return "", err
		}
		res = append(res, str)
	}

	return strings.Join(res, "\n"), nil
}
