package db

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/RobertsMJ/simc-cloud-backend/job"
	"github.com/RobertsMJ/simc-cloud-backend/models"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type jobItem struct {
	PK             string        `dynamodbav:"PK"`
	SK             string        `dynamodbav:"SK"`
	ID             string        `dynamodbav:"id"`
	Status         models.Status `dynamodbav:"status"`
	TotalCount     int           `dynamodbav:"total_count"`
	CompletedCount int           `dynamodbav:"completed_count"`
	FailedCount    int           `dynamodbav:"failed_count"`
	CreatedAt      string        `dynamodbav:"created_at"`
}

func jobItemFromResponse(r models.Job) jobItem {
	return jobItem{
		PK:             r.ID,
		SK:             "JOB",
		ID:             r.ID,
		Status:         r.Status,
		TotalCount:     r.TotalCount,
		CompletedCount: r.CompletedCount,
		FailedCount:    r.FailedCount,
		CreatedAt:      r.CreatedAt,
	}
}

type jobRepository struct {
	client    *dynamodb.Client
	tableName string
}

type JobRepositoryConfig struct {
	Client    *dynamodb.Client
	TableName string
}

func NewJobRepository(config JobRepositoryConfig) *jobRepository {
	return &jobRepository{
		client:    config.Client,
		tableName: config.TableName,
	}
}

var _ job.Writer = (*jobRepository)(nil)

func (r *jobRepository) WriteResult(ctx context.Context, result models.SimResult) error {
	item := resultItemFromResult(result)

	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return err
	}

	var counterField string
	switch result.Status {
	case models.StatusCompleted:
		counterField = "completed_count"
	case models.StatusError:
		counterField = "failed_count"
	default:
		return fmt.Errorf("invalid result status: %s", result.Status)
	}
	updateExpr := fmt.Sprintf("ADD %s :inc", counterField)

	_, err = r.client.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{
		TransactItems: []types.TransactWriteItem{
			// Atomically increment the completed / failed count and update the status of the job
			{
				Update: &types.Update{
					TableName: aws.String(r.tableName),
					Key: map[string]types.AttributeValue{
						"PK": &types.AttributeValueMemberS{Value: item.PK},
						"SK": &types.AttributeValueMemberS{Value: "JOB"},
					},
					UpdateExpression: aws.String(updateExpr),
				},
			},
			// Put the sim result item
			{
				Put: &types.Put{
					TableName: aws.String(r.tableName),
					Item:      av,
				},
			},
		},
	})
	if err != nil {
		slog.Error("failed to write result", slog.Any("error", err))
		return err
	}

	// Fetch the updated item
	updatedJob, err := r.GetJob(ctx, result.JobID)
	if err != nil {
		slog.Error("failed to fetch updated job", slog.Any("error", err))
		return err
	}

	if updatedJob == nil {
		slog.Error("job not found after update", slog.Any("job_id", result.JobID))
		return fmt.Errorf("job not found after update")
	}

	if updatedJob.TotalCount > 0 && updatedJob.CompletedCount+updatedJob.FailedCount >= updatedJob.TotalCount {
		// All sims are completed, update the job status to completed
		var finalStatus models.Status
		if updatedJob.FailedCount > 0 {
			finalStatus = models.StatusError
		} else {
			finalStatus = models.StatusCompleted
		}
		_, err = r.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
			TableName: aws.String(r.tableName),
			Key: map[string]types.AttributeValue{
				"PK": &types.AttributeValueMemberS{Value: item.PK},
				"SK": &types.AttributeValueMemberS{Value: "JOB"},
			},
			UpdateExpression: aws.String("SET #s = :finalStatus"),
			ExpressionAttributeNames: map[string]string{
				"#s": "status",
			},
			ConditionExpression: aws.String("#s = :inProgress"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":finalStatus": &types.AttributeValueMemberS{Value: string(finalStatus)},
				":inProgress":  &types.AttributeValueMemberS{Value: string(models.StatusInProgress)},
			},
		})

		var ccfe *types.ConditionalCheckFailedException
		if errors.As(err, &ccfe) {
			return nil // another invocation already finalized the job
		}
		if err != nil {
			slog.Error("failed to finalize job status", slog.Any("error", err))
			return err
		}
	}

	return nil
}

func (r *jobRepository) GetJob(ctx context.Context, jobID string) (*models.Job, error) {
	resp, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: jobID},
			"SK": &types.AttributeValueMemberS{Value: "JOB"},
		},
	})
	if err != nil {
		return nil, err
	}

	if resp.Item == nil {
		return nil, nil
	}

	var item jobItem
	err = attributevalue.UnmarshalMap(resp.Item, &item)
	if err != nil {
		return nil, err
	}

	return &models.Job{
		ID:             item.ID,
		Status:         item.Status,
		TotalCount:     item.TotalCount,
		CompletedCount: item.CompletedCount,
		FailedCount:    item.FailedCount,
		CreatedAt:      item.CreatedAt,
	}, nil
}
