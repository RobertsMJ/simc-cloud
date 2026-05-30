package models

type CreateJobRequest struct{}

type CreateJobResponse struct{}

type Job struct {
	ID string `json:"id" dynamodbav:"id"`
}
