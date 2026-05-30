package db

import (
	"context"

	"github.com/RobertsMJ/simc-cloud-backend/models"
)

type simResultItem struct {
	PK           string               `dynamodbav:"PK"`
	SK           string               `dynamodbav:"SK"`
	JobID        string               `dynamodbav:"job_id"`
	GearsetID    string               `dynamodbav:"gearset_id"`
	Status       models.SimStatus     `dynamodbav:"status"`
	ErrorMessage *string              `dynamodbav:"error_message,omitempty"`
	Baseline     bool                 `dynamodbav:"baseline,omitempty"`
	Statistics   models.SimStatistics `dynamodbav:"statistics"`
	Metadata     *map[string]any      `dynamodbav:"metadata,omitempty"`
	Result       *models.SimcOutput   `dynamodbav:"result,omitempty"`
}

func itemFromResponse(r models.SimResult) simResultItem {
	return simResultItem{
		PK:           r.JobID,
		SK:           "RESULT#" + r.GearsetID,
		JobID:        r.JobID,
		GearsetID:    r.GearsetID,
		Status:       r.Status,
		ErrorMessage: r.ErrorMessage,
		Baseline:     r.Baseline,
		Statistics:   r.Statistics,
		Metadata:     r.Metadata,
		Result:       r.Result,
	}
}

// Publishes the sim result to an SQS queue for processing by the database writer lambda.
// This decouples the sim runner from the database and allows for better scalability and
// fault tolerance. The sim runner can publish the result to the queue and then return
// immediately, while the database writer lambda can process the results from the queue
// at its own pace.
type simRepository struct {
}

func NewSimRepository() *simRepository {
	return &simRepository{}
}

func (r *simRepository) PublishResult(ctx context.Context, result models.SimResult) error {
	return models.ErrNotImplemented
}
