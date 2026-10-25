package db

type ResultRepository struct{}

func NewResultRepository() *ResultRepository {
	return &ResultRepository{}
}

// TODO: contract and implement for results
func (r *ResultRepository) Create() error {
	panic("not implemented")
}
