package job

import "github.com/RobertsMJ/simc-cloud-backend/models"

type CreateJobInput struct{}
type CreateJobOutput struct{}

type Service interface {
	CreateJob(input CreateJobInput) (CreateJobOutput, error)
}

type service struct{}

func NewService() Service {
	return &service{}
}

func (s *service) CreateJob(input CreateJobInput) (CreateJobOutput, error) {
	return CreateJobOutput{}, models.ErrNotImplemented
}
