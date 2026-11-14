package simc

import (
	"strings"

	"github.com/RobertsMJ/simc-cloud/backend/simc/internal"
)

type Profile struct {
	Character Character      `json:"character"`
	Baseline  Loadout        `json:"loadout"`
	Options   LoadoutOptions `json:"options"`
	Config    Config         `json:"config"`
	Extra     []Parameter    `json:"extra,omitempty"`
}

var (
	_ internal.ValueMarshaler = Profile{}
)

func (p Profile) MarshalSimc() (string, error) {
	character, err := p.Character.MarshalSimc()
	if err != nil {
		return "", err
	}

	baseline, err := p.Baseline.MarshalSimc()
	if err != nil {
		return "", err
	}

	options, err := p.Options.MarshalSimc()
	if err != nil {
		return "", err
	}

	config, err := p.Config.MarshalSimc()
	if err != nil {
		return "", err
	}

	return strings.Join([]string{character, baseline, options, config}, "\n"), nil
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

	if statement.Disabled {
		if p.Options.Items == nil {
			p.Options.Items = make(map[SlotGroup][]Item)
		}
		if p.Options.Items[item.SlotGroup] == nil {
			p.Options.Items[item.SlotGroup] = make([]Item, 0)
		}
		p.Options.Items[item.SlotGroup] = append(p.Options.Items[item.SlotGroup], item)
	} else {
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

	if statement.Disabled {
		p.Options.OmniumTalents = append(p.Options.OmniumTalents, om)
	} else {
		p.Baseline.OmniumTalents = om
	}
	return nil
}
