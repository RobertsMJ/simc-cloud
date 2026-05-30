package models

// Models for single simulations
// Request: Single job and gearset with a simc input string
// Before the sim is run, the job and gearset will be persisted in the database with a status of "in_progress". After the sim is run, the result will be persisted in the database with a status of "completed" or "error" depending on the outcome of the sim. The response will include the job ID, gearset ID, status, and result (if successful) or error message (if failed).
// Response: Job, Gearset, and simc json output to persist in database and return to client

type SimRequest struct {
	JobID     string          `json:"job_id"`
	GearsetID string          `json:"gearset_id"`
	Metadata  *map[string]any `json:"metadata,omitempty"`
	Input     string          `json:"input"`
}

type SimStatus string

const (
	StatusInProgress SimStatus = "in_progress"
	StatusCompleted  SimStatus = "completed"
	StatusError      SimStatus = "error"
)

type SimcOutput map[string]any

type StatSample struct {
	Sum          float64 `json:"sum"`
	Count        int     `json:"count"`
	Mean         float64 `json:"mean"`
	Min          float64 `json:"min,omitempty"`
	Max          float64 `json:"max,omitempty"`
	Median       float64 `json:"median,omitempty"`
	Variance     float64 `json:"variance,omitempty"`
	StdDev       float64 `json:"std_dev,omitempty"`
	MeanVariance float64 `json:"mean_variance,omitempty"`
	MeanStdDev   float64 `json:"mean_std_dev,omitempty"`
}

type SimStatistics struct {
	ElapsedCpuSeconds    float64    `json:"elapsed_cpu_seconds"`
	ElapsedTimeSeconds   float64    `json:"elapsed_time_seconds"`
	InitTimeSeconds      float64    `json:"init_time_seconds"`
	MergeTimeSeconds     float64    `json:"merge_time_seconds"`
	AnalyzeTimeSeconds   float64    `json:"analyze_time_seconds"`
	TotalEventsProcessed int64      `json:"total_events_processed"`
	SimulationLength     StatSample `json:"simulation_length"`
	RaidDps              StatSample `json:"raid_dps"`
	TotalDmg             StatSample `json:"total_dmg"`
}

type SimResult struct {
	JobID        string          `json:"job_id"`
	GearsetID    string          `json:"gearset_id"`
	Status       SimStatus       `json:"status"`
	ErrorMessage *string         `json:"error_message,omitempty"`
	Baseline     bool            `json:"baseline"`
	Statistics   SimStatistics   `json:"statistics"`
	Metadata     *map[string]any `json:"metadata,omitempty"`
	Result       *SimcOutput     `json:"result"`
}
