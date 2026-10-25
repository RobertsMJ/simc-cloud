package db

type JobRepository struct{}

func NewJobRepository() *JobRepository {
	return &JobRepository{}
}

// TODO: contract and implement for results
func (r *JobRepository) Create() error {
	panic("not implemented")
}
