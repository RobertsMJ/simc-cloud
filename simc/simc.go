package simc

import (
	"io"
	"log/slog"

	"github.com/RobertsMJ/simc-cloud/simc/internal"
)

// Some expected strings:
// waist=,id=244573,bonus_id=12214/13667/12497/12066/8960/12384/8791/13622/12667,content_tuning=3615,crafted_stats=49/32,crafting_quality=5
// server=illidan
// omnium_talents=136819:1/136822:1

type SimulationDocument struct {
	Character      Character
	ActiveLoadout  Loadout
	LoadoutOptions LoadoutOptions
	Extra          []Parameter
}

func NewSimulationDocument(doc io.Reader) (sd SimulationDocument, err error) {
	statements, err := internal.ParseStatements(doc)
	if err != nil {
		return SimulationDocument{}, err
	}

	sd.Character = Character{}

	for _, statement := range statements {
		switch {
		case isCharacterKey(statement.Key):
			if err := sd.Character.UnmarshalStatement(statement); err != nil {
				slog.Error("could not unmarshal character statement", "err", err, "statement", statement)
				return sd, err
			}
		case statement.Key == "talents":
			if err = sd.handleTalent(statement); err != nil {
				slog.Error("could not handle talent statement", "err", err, "statement", statement)
				return sd, err
			}
		case isItemKey(statement.Key):
			if err = sd.handleItem(statement); err != nil {
				slog.Error("could not handle item statement", "err", err, "statement", statement)
				return sd, err
			}
		case statement.Key == "omnium_talents":
			if err = sd.handleOmniumTalents(statement); err != nil {
				slog.Error("could not handle omnium talents", "err", err, "statement", statement)
				return sd, err
			}
		default:
			if sd.Extra == nil {
				sd.Extra = []Parameter{}
			}
			sd.Extra = append(sd.Extra, Parameter{
				Key:   statement.Key,
				Value: statement.Value,
			})
		}
	}

	return sd, nil
}

func (sd *SimulationDocument) handleTalent(statement internal.Statement) error {
	var talent Talent
	if err := talent.UnmarshalStatement(statement); err != nil {
		return err
	}
	if statement.Disabled {
		if sd.LoadoutOptions.Talents == nil {
			sd.LoadoutOptions.Talents = []Talent{}
		}
		sd.LoadoutOptions.Talents = append(sd.LoadoutOptions.Talents, talent)
		return nil
	} else {
		sd.ActiveLoadout.Talent = talent
		return nil
	}
}

func (sd *SimulationDocument) handleItem(statement internal.Statement) error {
	var item Item
	if err := item.UnmarshalStatement(statement); err != nil {
		return err
	}

	if statement.Disabled {
		if sd.LoadoutOptions.Items == nil {
			sd.LoadoutOptions.Items = make([]Item, 0, 1)
		}
		sd.LoadoutOptions.Items = append(sd.LoadoutOptions.Items, item)
	} else {
		if sd.ActiveLoadout.Items == nil {
			sd.ActiveLoadout.Items = make(map[ItemSlot]Item)
		}
		sd.ActiveLoadout.Items[item.Slot] = item
	}
	return nil
}

func (sd *SimulationDocument) handleOmniumTalents(statement internal.Statement) error {
	om := internal.IDValueList{}
	if err := om.UnmarshalStatement(statement); err != nil {
		return err
	}

	if statement.Disabled {
		sd.LoadoutOptions.OmniumTalents = append(sd.LoadoutOptions.OmniumTalents, om)
	} else {
		sd.ActiveLoadout.OmniumTalents = om
	}
	return nil
}

func Marshal(doc *SimulationDocument) (io.Reader, error) {
	panic("not implemented")
}
