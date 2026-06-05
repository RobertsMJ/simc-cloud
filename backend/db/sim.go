package db

import (
	"github.com/RobertsMJ/simc-cloud-backend/models"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type resultItem struct {
	PK           string               `dynamodbav:"PK"`
	SK           string               `dynamodbav:"SK"`
	JobID        string               `dynamodbav:"job_id"`
	GearsetID    string               `dynamodbav:"gearset_id"`
	Status       models.Status        `dynamodbav:"status"`
	ErrorMessage *string              `dynamodbav:"error_message,omitempty"`
	IsBaseline   *bool                `dynamodbav:"is_baseline,omitempty"`
	Statistics   models.SimStatistics `dynamodbav:"statistics"`
	Metadata     *map[string]any      `dynamodbav:"metadata,omitempty"`
	Result       []byte               `dynamodbav:"result,omitempty"`
}

func resultItemFromResult(r models.SimResult) resultItem {
	return resultItem{
		PK:           r.JobID,
		SK:           "RESULT#" + r.GearsetID,
		JobID:        r.JobID,
		GearsetID:    r.GearsetID,
		Status:       r.Status,
		ErrorMessage: r.ErrorMessage,
		IsBaseline:   &r.Baseline,
		Statistics:   r.Statistics,
		Metadata:     r.Metadata,
		Result:       r.Result,
	}
}

type simRepository struct {
	client    *dynamodb.Client
	tableName string
}

type SimRepositoryConfig struct {
	Client    *dynamodb.Client
	TableName string
}

func NewSimRepository(config SimRepositoryConfig) *simRepository {
	return &simRepository{
		client:    config.Client,
		tableName: config.TableName,
	}
}

// TODO:MJR - Should only handle fetching results
