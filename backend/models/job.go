package models

import "github.com/RobertsMJ/simc-cloud/simc"

type CreateJobRequest struct {
	Character simc.Character      `json:"character"`
	Baseline  simc.Loadout        `json:"baseline"`
	Options   simc.LoadoutOptions `json:"options"`
}

type CreateJobResponse struct {
	ID     string `json:"id"`
	Status Status `json:"status"`
}

type Status string

const (
	StatusInProgress Status = "in_progress"
	StatusCompleted  Status = "completed"
	StatusError      Status = "error"
)

type Job struct {
	ID             string `json:"id"`
	Status         Status `json:"status"`
	TotalCount     int    `json:"total_count"`
	CompletedCount int    `json:"completed_count"`
	FailedCount    int    `json:"failed_count"`
	CreatedAt      string `json:"created_at"`
}
