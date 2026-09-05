package job

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
	panic("not implemented")
}
