package simc

import (
	"fmt"
	"maps"
	"slices"
	"strings"
)

// Loadout represents a character's active gear and talent configuration.
type Loadout struct {
	Items         map[ItemSlot]Item `json:"items"`
	Talent        Talent            `json:"talents"`
	OmniumTalents OmniumTalents     `json:"omnium_talents"`
}

func (l Loadout) render() string {
	res := []string{}

	for _, slot := range slices.Sorted(maps.Keys(l.Items)) {
		res = append(res, l.Items[slot].render(slot))
	}

	return strings.Join(res, "\n")
}

func (l Loadout) renderWithID(id string) string {
	res := []string{}

	for i, slot := range slices.Sorted(maps.Keys(l.Items)) {
		if i == 0 {
			res = append(res, fmt.Sprintf(`profileset."%s"=%s`, id, l.Items[slot].render(slot)))
		} else {
			res = append(res, fmt.Sprintf(`profileset."%s"+=%s`, id, l.Items[slot].render(slot)))
		}
	}

	if len(l.Talent.Value) > 0 {
		if len(res) > 0 {
			res = append(res, fmt.Sprintf(`profileset."%s"+=%s`, id, l.Talent.render()))
		} else {
			res = append(res, fmt.Sprintf(`profileset."%s"=%s`, id, l.Talent.render()))
		}
	}

	if len(l.OmniumTalents.Value) > 0 {
		if len(res) > 0 {
			res = append(res, fmt.Sprintf(`profileset."%s"+=%s`, id, l.OmniumTalents.render()))
		} else {
			res = append(res, fmt.Sprintf(`profileset."%s"=%s`, id, l.OmniumTalents.render()))
		}
	}

	return strings.Join(res, "\n")
}
