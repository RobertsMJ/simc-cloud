package models

type CreateJobRequest struct{}

type CreateJobResponse struct{}

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
