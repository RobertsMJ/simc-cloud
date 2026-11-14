package models

import (
	"github.com/RobertsMJ/simc-cloud/backend/simc"
)

type CreateJobRequestBody struct {
	Character simc.Character       `json:"character" validate:"required"`
	Baseline  simc.Loadout         `json:"baseline" validate:"required"`
	Options   *simc.LoadoutOptions `json:"options,omitempty"`
	Config    simc.Config          `json:"config"`
}

type CreateJobResponseBody struct {
	Id     string `json:"id"`
	Status Status `json:"status"`
}

type Status string

const (
	StatusPending    Status = "pending"
	StatusInProgress Status = "in_progress"
	StatusDone       Status = "done"
	StatusError      Status = "error"
)

type Job struct {
	ID        string `json:"id"`
	Status    Status `json:"status"`
	CreatedAt string `json:"created_at"`
}
