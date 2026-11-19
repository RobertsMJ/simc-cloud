package models

import "github.com/RobertsMJ/simc-cloud/backend/simc"

type Gearset struct {
	ID    string                      `json:"id,omitempty"`
	Items map[simc.ItemSlot]simc.Item `json:"items"`

	Talent        simc.Talent      `json:"talents"`
	OmniumTalents simc.IDValueList `json:"omnium_talents,omitempty"`
}

func (sl Gearset) ToSimcLoadout() simc.Loadout {
	return simc.Loadout{
		ID:            sl.ID,
		Items:         sl.Items,
		Talent:        sl.Talent,
		OmniumTalents: sl.OmniumTalents,
	}
}

type GearsetOptions struct {
	Items   []simc.Item          `json:"items"`
	Talents []simc.Talent        `json:"talents"`
	Omnium  []simc.OmniumTalents `json:"omnium"`
}

type SimcJobConfig struct {
	Character simc.Character `json:"character"`
	Baseline  Gearset        `json:"baseline"`
	Options   GearsetOptions `json:"options"`
	Config    simc.Config    `json:"config"`
}
