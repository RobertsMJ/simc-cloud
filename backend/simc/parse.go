package simc

import (
	"fmt"
	"io"
)

type AddonExport struct {
	Character    Character         `json:"character"`
	Equipped     Loadout           `json:"equipped"`
	Bags         []Item            `json:"bags"`
	SavedTalents []Talent          `json:"savedTalents"`
	Metadata     map[string]string `json:"metadata"`
}

// Parse parses the provided reader into an AddonExport structure.
func Parse(r io.Reader) (AddonExport, error) {
	statements, err := parseStatements(r)
	if err != nil {
		return AddonExport{}, err
	}

	characterStatements := []statement{}
	bagItems := []Item{}
	activeLoadout := Loadout{
		Items: map[ItemSlot]Item{},
	}
	savedTalents := []Talent{}
	metadata := map[string]string{}
	for _, s := range statements {
		switch {
		case isCharacterKey(s.Key):
			characterStatements = append(characterStatements, s)
		case isItemKey(s.Key):
			item, err := newItem(s)
			if err != nil {
				return AddonExport{}, ParseError{Line: s.Line, Err: err}
			}

			if s.Disabled {
				bagItems = append(bagItems, item)
			} else {
				for _, slot := range item.Slot.Slots() {
					if activeLoadout.Items[slot].ID == 0 {
						activeLoadout.Items[slot] = item
						break
					}
				}
			}
		case s.Key == talentKey:
			talent, err := newTalent(s.Comment, s.Value)
			if err != nil {
				return AddonExport{}, ParseError{Line: s.Line, Err: err}
			}
			if s.Disabled {
				savedTalents = append(savedTalents, talent)
			} else {
				activeLoadout.Talents = talent
			}
		case s.Key == omniumTalentKey:
			om, err := newOmniumTalents(s.Value)
			if err != nil {
				return AddonExport{}, ParseError{Line: s.Line, Err: err}
			}
			activeLoadout.OmniumTalents = om
		default:
			metadata[s.Key] = s.Value
		}
	}

	character, err := newCharacter(characterStatements)
	if err != nil {
		return AddonExport{}, ParseError{Line: 0, Err: err}
	}

	return AddonExport{
		Character:    character,
		Equipped:     activeLoadout,
		Bags:         bagItems,
		SavedTalents: savedTalents,
		Metadata:     metadata,
	}, nil
}

type ParseError struct {
	Line int
	Err  error
}

func (e ParseError) Error() string {
	return fmt.Sprintf("line %d: %v", e.Line, e.Err)
}

func (e ParseError) Unwrap() error {
	return e.Err
}
