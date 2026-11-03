package natsx

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/suite"
)

type MetricsTestSuite struct {
	suite.Suite
}

func TestMetricsTestSuite(t *testing.T) {
	suite.Run(t, new(MetricsTestSuite))
}

func (s *MetricsTestSuite) TestConsumerMetrics() {
	metrics := NewMetrics(prometheus.NewRegistry())
	metrics.forConsumer(&ConsumerConfig{StreamName: "stream", DurableName: "consumer"})
}

func (s *MetricsTestSuite) TestPublisherMetrics() {
	metrics := NewMetrics(prometheus.NewRegistry())
	metrics.forPublisher(&PublisherConfig{StreamName: "stream"})
}
