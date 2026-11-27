package models

import (
	"github.com/RobertsMJ/simc-cloud/backend/simc"
)

// Models for single simulations
// Request: Single job and gearset with a simc input string
// Before the sim is run, the job and gearset will be persisted in the database with a status of "in_progress". After the sim is run, the result will be persisted in the database with a status of "completed" or "error" depending on the outcome of the sim. The response will include the job ID, gearset ID, status, and result (if successful) or error message (if failed).
// Response: Job, Gearset, and simc json output to persist in database and return to client
type RunSim struct {
	JobID    string                  `json:"job_id"`
	Gearsets map[string]simc.Loadout `json:"gearsets"`
	Input    string                  `json:"input"`
}

type SimCompleted struct {
	JobID        string                  `json:"job_id"`
	Status       Status                  `json:"status"`
	ErrorMessage *string                 `json:"error_message,omitempty"`
	Gearsets     map[string]simc.Loadout `json:"gearsets"`
	Result       simc.Report             `json:"result"`
}
