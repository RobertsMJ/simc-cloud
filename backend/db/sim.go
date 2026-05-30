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

type simRepository struct {
	// TODO: DynamoDB client and table name
}

func NewSimRepository() *simRepository {
	return &simRepository{}
}

func (r *simRepository) SaveSimResult(ctx context.Context, result models.SimResult) error {
	// TODO: Save sim result to DynamoDB using itemFromResponse to convert to simResultItem
	return models.ErrNotImplemented
}
