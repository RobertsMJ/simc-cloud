package simc

import (
	"io"
	"log/slog"

	"github.com/RobertsMJ/simc-cloud/backend/simc/internal"
)

// Profile represents a parsed simc addon string with character, active loadout, options (for items, talents), config, and a catch-all forextra parameters.
// Need to be able to go from a simc string to this representation
type Profile struct {
	Character Character      `json:"character"`
	Baseline  Loadout        `json:"loadout"`
	Options   LoadoutOptions `json:"options"`
	Config    Config         `json:"config"`
	Extra     []Parameter    `json:"extra,omitempty"`
}

func NewProfile(doc io.Reader) (*Profile, error) {
	p := &Profile{}

	statements, err := internal.ParseStatements(doc)
	if err != nil {
		return nil, err
	}

	for _, statement := range statements {
		switch {
		case isCharacterKey(statement.Key):
			if err := p.Character.UnmarshalStatement(statement); err != nil {
				slog.Error("could not unmarshal character statement", "err", err, "statement", statement)
				return nil, err
			}
		case statement.Key == "talents":
			if err = p.unmarshalTalent(statement); err != nil {
				slog.Error("could not handle talent statement", "err", err, "statement", statement)
				return nil, err
			}
		case isItemKey(statement.Key):
			if err = p.unmarshalItem(statement); err != nil {
				slog.Error("could not handle item statement", "err", err, "statement", statement)
				return nil, err
			}
		case statement.Key == "omnium_talents":
			if err = p.unmarshalOmniumTalents(statement); err != nil {
				slog.Error("could not handle omnium talents", "err", err, "statement", statement)
				return nil, err
			}
		default:
			if p.Extra == nil {
				p.Extra = []Parameter{}
			}
			p.Extra = append(p.Extra, Parameter{
				Key:   statement.Key,
				Value: statement.Value,
			})
		}
	}

	return p, nil
}

func (p *Profile) unmarshalTalent(statement internal.Statement) error {
	var talent Talent
	if err := talent.UnmarshalStatement(statement); err != nil {
		return err
	}
	if statement.Disabled {
		if p.Options.Talents == nil {
			p.Options.Talents = []Talent{}
		}
		p.Options.Talents = append(p.Options.Talents, talent)
		return nil
	} else {
		p.Baseline.Talent = talent
		return nil
	}
}

func (p *Profile) unmarshalItem(statement internal.Statement) error {
	var item Item
	if err := item.UnmarshalStatement(statement); err != nil {
		return err
	}

	// Disabled items represent items in the bag
	// They are stored in the options under the corresponding slot group.
	if statement.Disabled {
		if p.Options.Items == nil {
			p.Options.Items = make(map[SlotGroup][]Item)
		}
		if p.Options.Items[item.SlotGroup] == nil {
			p.Options.Items[item.SlotGroup] = make([]Item, 0)
		}
		p.Options.Items[item.SlotGroup] = append(p.Options.Items[item.SlotGroup], item)
	} else {
		// Otherwise it's an equipped item
		if p.Baseline.Items == nil {
			p.Baseline.Items = make(map[ItemSlot]Item)
		}
		p.Baseline.Items[item.Slot] = item
	}

	return nil
}

func (p *Profile) unmarshalOmniumTalents(statement internal.Statement) error {
	om := IDValueList{}
	if err := om.UnmarshalStatement(statement); err != nil {
		return err
	}

	// Addon is only going to report the active omnium talents
	if !statement.Disabled {
		p.Baseline.OmniumTalents = om
	}
	return nil
}
