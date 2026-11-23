package simc

import (
	"crypto/sha256"
	"fmt"
	"maps"
	"slices"
	"strings"
)

// Loadout represents a character's active gear and talent configuration.
type Loadout struct {
	Items         map[ItemSlot]Item `json:"items"`
	Talents       Talent            `json:"talents"`
	OmniumTalents OmniumTalents     `json:"omnium_talents"`
}

func (l Loadout) Hash() string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(l.render())))
}

func (l Loadout) getStrings() []string {
	res := []string{}
	for _, slot := range slices.Sorted(maps.Keys(l.Items)) {
		res = append(res, l.Items[slot].render(slot))
	}
	res = append(res, l.Talents.render(), l.OmniumTalents.render("omnium_talents"))
	return res
}

func (l Loadout) render() string {
	return strings.Join(l.getStrings(), "\n")
}

func (l Loadout) renderWithID(id string) string {
	str := l.getStrings()
	for i, s := range str {
		if i == 0 {
			str[i] = fmt.Sprintf(`profileset."%s"=%s`, id, s)
		} else {
			str[i] = fmt.Sprintf(`profileset."%s"+=%s`, id, s)
		}
	}

	return strings.Join(str, "\n")
}
