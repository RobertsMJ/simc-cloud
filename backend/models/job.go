package models

import (
	"github.com/RobertsMJ/simc-cloud/backend/simc"
)

type Job struct {
	ID        string `json:"id"`
	Status    Status `json:"status"`
	CreatedAt string `json:"created_at"`
}

type Status string

const (
	StatusPending    Status = "pending"
	StatusInProgress Status = "in_progress"
	StatusDone       Status = "done"
	StatusError      Status = "error"
)

type JobSpec struct {
	Character simc.Character `json:"character"`
	Base      simc.Loadout   `json:"base"`
	Options   Candidates     `json:"options"`
	Config    simc.Config    `json:"config"`
}
type Candidates struct {
	Items   []simc.Item          `json:"items"`
	Talents []simc.Talent        `json:"talents"`
	Omnium  []simc.OmniumTalents `json:"omnium"`
}
