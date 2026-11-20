package simc

import (
	"io"

	"github.com/RobertsMJ/simc-cloud/backend/simc/internal"
)

type AddonExport struct {
	Character    Character `json:"character"`
	Equipped     Loadout   `json:"equipped"`
	Bags         []Item    `json:"bags"`
	SavedTalents []Talent  `json:"savedTalents"`
}

func ParseAddonExport(r io.Reader) (AddonExport, error) {
	statements, err := internal.ParseStatements(r)
	if err != nil {
		return AddonExport{}, err
	}

	characterStatements := []internal.Statement{}
	bagItems := []Item{}
	activeLoadout := Loadout{
		Items: map[ItemSlot]Item{},
	}
	savedTalents := []Talent{}
	for _, s := range statements {
		if isCharacterStatement(s) {
			characterStatements = append(characterStatements, s)
			continue
		}

		if isItemStatement(s) {
			item, err := newItem(s)
			if err != nil {
				return AddonExport{}, err
			}

			if s.Disabled {
				bagItems = append(bagItems, item)
			} else {
				switch item.Slot {
				case SlotGroupFinger:
					if activeLoadout.Items[ItemSlotFinger1].ID == 0 {
						activeLoadout.Items[ItemSlotFinger1] = item
					} else {
						activeLoadout.Items[ItemSlotFinger2] = item
					}
				case SlotGroupTrinket:
					if activeLoadout.Items[ItemSlotTrinket1].ID == 0 {
						activeLoadout.Items[ItemSlotTrinket1] = item
					} else {
						activeLoadout.Items[ItemSlotTrinket2] = item
					}
				default:
					activeLoadout.Items[ItemSlot(item.Slot)] = item
				}
			}
			continue
		}

		if s.Key == talentKey {
			talent, err := newTalent(s)
			if err != nil {
				return AddonExport{}, err
			}
			if s.Disabled {
				savedTalents = append(savedTalents, talent)
			} else {
				activeLoadout.Talent = talent
			}
			continue
		}

		if s.Key == omniumTalentKey {
			om, err := newOmniumTalents(s)
			if err != nil {
				return AddonExport{}, err
			}
			activeLoadout.OmniumTalents = om
		}
	}
	character, err := NewCharacter(characterStatements)
	if err != nil {
		return AddonExport{}, err
	}

	return AddonExport{
		Character:    character,
		Equipped:     activeLoadout,
		Bags:         bagItems,
		SavedTalents: savedTalents,
	}, nil
}

func isCharacterStatement(s internal.Statement) bool {
	if _, err := NewWowClass(s.Key); err == nil {
		return true
	}
	switch s.Key {
	case charFieldLevel,
		charFieldRole,
		charFieldSpec,
		charFieldRace,
		charFieldRegion,
		charFieldServer:
		return true
	}
	return false
}

func isItemStatement(s internal.Statement) bool {
	if _, err := NewItemSlot(s.Key); err == nil {
		return true
	}
	return false
}
