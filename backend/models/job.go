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
	Character simc.Character
	Base      simc.Loadout
	Options   Candidates
	Config    simc.Config
}
type Candidates struct {
	Items   []simc.Item
	Talents []simc.Talent
	Omnium  []simc.OmniumTalents
}
