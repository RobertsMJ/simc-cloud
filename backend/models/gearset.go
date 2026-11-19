package models

import "github.com/RobertsMJ/simc-cloud/backend/simc"

type GenerateGearsetRequest struct {
	JobID     string         `json:"job_id"`
	Character simc.Character `json:"character"`
	Baseline  Gearset        `json:"baseline"`
	Options   GearsetOptions `json:"options"`
	Config    simc.Config    `json:"config"`
}
