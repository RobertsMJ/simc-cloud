package db_test

import (
	"context"
	"testing"

	"github.com/RobertsMJ/simc-cloud-backend/config"
	"github.com/RobertsMJ/simc-cloud-backend/db"
	"github.com/RobertsMJ/simc-cloud-backend/models"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/stretchr/testify/suite"
)

type JobRepositoryTestSuite struct {
	suite.Suite
	ctx       context.Context
	client    *dynamodb.Client
	tableName string

	jobRepo db.JobRepository
}

func TestJobRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(JobRepositoryTestSuite))
}

func (suite *JobRepositoryTestSuite) SetupTest() {
	suite.ctx = context.Background()
	cfg := config.LoadAWS(context.Background())
	cfg.BaseEndpoint = aws.String("http://localhost:8000")

	suite.client = dynamodb.NewFromConfig(cfg)
	suite.tableName = "sim-table"

	suite.jobRepo = db.NewJobRepository(db.JobRepositoryConfig{
		Client:    suite.client,
		TableName: suite.tableName,
	})
}

func (suite *JobRepositoryTestSuite) TestCreateJob() {
	job := models.Job{
		ID:             "job-123",
		Status:         models.StatusInProgress,
		TotalCount:     10,
		CompletedCount: 0,
		FailedCount:    0,
		CreatedAt:      "2024-06-01T12:00:00Z",
	}

	err := suite.jobRepo.CreateJob(suite.ctx, job)
	suite.NoError(err)

	retrievedJob, err := suite.jobRepo.GetJob(suite.ctx, job.ID)
	suite.NoError(err)
	suite.Equal(job, retrievedJob)
}
