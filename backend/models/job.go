package models

import "github.com/RobertsMJ/simc-cloud-backend/simc"

type CreateJobRequest struct {
	Character simc.Character
	Baseline  simc.Loadout
	Options   simc.LoadoutOptions
	Config    simc.SimConfig // TODO:MJR How should I handle this?
}

type CreateJobResponse struct {
	ID     string         `json:"id"`
	Status Status         `json:"status"`
	Config simc.SimConfig `json:"config"`
}

type Status string

const (
	StatusInProgress Status = "in_progress"
	StatusCompleted  Status = "completed"
	StatusError      Status = "error"
)

type Job struct {
	ID             string         `json:"id"`
	Status         Status         `json:"status"`
	Config         simc.SimConfig `json:"config"`
	TotalCount     int            `json:"total_count"`
	CompletedCount int            `json:"completed_count"`
	FailedCount    int            `json:"failed_count"`
	CreatedAt      string         `json:"created_at"`
}
