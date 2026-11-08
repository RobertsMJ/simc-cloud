package sim

import (
	"context"
	"testing"
	"time"

	"github.com/RobertsMJ/simc-cloud/backend/models"
	"github.com/stretchr/testify/suite"
)

type ServiceTestSuite struct {
	suite.Suite

	publisher *resultPublisher
}

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}

func (s *ServiceTestSuite) SetupTest() {
	s.publisher = &resultPublisher{}
}

type simulatorFunc func(context.Context, Input) (Output, error)

func (f simulatorFunc) Run(ctx context.Context, in Input) (Output, error) {
	return f(ctx, in)
}

type resultPublisher struct {
	results []models.SimResult
}

func (p *resultPublisher) Publish(_ context.Context, result models.SimResult) error {
	p.results = append(p.results, result)
	return nil
}

func (s *ServiceTestSuite) TestHandle_CallerCancellationIsReturned() {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	service := NewService(
		simulatorFunc(func(context.Context, Input) (Output, error) {
			return Output{}, nil
		}),
		s.publisher,
		&Config{MaxRuntime: time.Second},
	)

	err := service.Handle(ctx, models.SimRequest{})

	s.ErrorIs(err, context.Canceled)
	s.Empty(s.publisher.results)
}
