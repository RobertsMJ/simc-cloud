package models

type SimulationRequest struct {
	RequestID string          `json:"request_id"`
	GearsetID string          `json:"gearset_id"`
	Metadata  *map[string]any `json:"metadata,omitempty"`
	Input     string          `json:"input"`
}

type SimulationResponse struct {
	RequestID string          `json:"request_id"`
	GearsetID string          `json:"gearset_id"`
	Metadata  *map[string]any `json:"metadata,omitempty"`
	Result    string          `json:"result"`
}
