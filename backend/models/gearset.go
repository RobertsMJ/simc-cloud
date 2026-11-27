package models

import "github.com/RobertsMJ/simc-cloud/backend/simc"

type GenerateGearsets struct {
	JobID string
	Spec  JobSpec
}

// GearsetResult is one profileset's result.
// ID is the profileset name, simc.Loadout.Hash().
type GearsetResult struct {
	JobID     string       `json:"job_id" db:"job_id"`
	ID        string       `json:"gearset_id" db:"gearset_id"`
	Loadout   simc.Loadout `json:"loadout" db:"loadout"`
	DPS       float64      `json:"dps" db:"dps"`
	DPSStdErr float64      `json:"dps_std_err" db:"dps_std_err"`
}
