package db

import (
	"context"

	"github.com/RobertsMJ/simc-cloud/backend/models"
)

type JobRepository struct{}

func NewJobRepository() *JobRepository {
	return &JobRepository{}
}

// TODO: contract and implement for results
func (r *JobRepository) Create(context.Context, models.Job) error {
	return nil
}
