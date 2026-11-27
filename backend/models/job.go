package models

import (
	"time"

	"github.com/RobertsMJ/simc-cloud/backend/simc"
)

type Job struct {
	ID                string     `json:"id" db:"id"`
	Status            Status     `json:"status" db:"status"`
	Spec              JobSpec    `json:"spec" db:"spec"`
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`
	CompletedAt       *time.Time `json:"completed_at,omitempty" db:"completed_at"`
	SimcVersion       *string    `json:"simc_version,omitempty" db:"simc_version"`
	BaselineDPS       *float64   `json:"baseline_dps,omitempty" db:"baseline_dps"`
	BaselineDPSStdErr *float64   `json:"baseline_dps_std_err,omitempty" db:"baseline_dps_std_err"`
	Error             *string    `json:"error,omitempty" db:"error"`
}

// Status values are constrained by the jobs.status CHECK constraint.
type Status string

const (
	StatusPending Status = "pending"
	StatusDone    Status = "done"
	StatusError   Status = "error"
)

type JobSpec struct {
	Character simc.Character `json:"character"`
	Base      simc.Loadout   `json:"base"`
	Options   Candidates     `json:"options"`
	Config    simc.Config    `json:"config"`
}

type Candidates struct {
	Items         []simc.Item          `json:"items"`
	Talents       []simc.Talent        `json:"talents"`
	OmniumTalents []simc.OmniumTalents `json:"omnium_talents"`
}

type JobResult struct {
	JobID             string
	SimcVersion       string
	BaselineDPS       float64
	BaselineDPSStdErr float64
	Gearsets          []GearsetResult
}
