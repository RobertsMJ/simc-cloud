package natsx_test

import (
	"testing"

	"github.com/RobertsMJ/simc-cloud/backend/platform/natsx"
	"github.com/stretchr/testify/suite"
)

type ConfigTestSuite struct {
	suite.Suite
}

func TestConfigTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}

// These keys are emitted by the Helm chart (_consumer.tpl, _publisher.tpl, values default env).
func (s *ConfigTestSuite) TestReadsChartEnvKeys() {
	s.T().Setenv("QUEUE_URL", "nats://nats:4222")
	s.T().Setenv("CONSUMER_STREAM_NAME", "in-stream")
	s.T().Setenv("CONSUMER_NAME", "in-consumer")
	s.T().Setenv("PUBLISHER_STREAM_NAME", "out-stream")
	s.T().Setenv("PUBLISHER_SUBJECT", "out.subject")

	nc, err := natsx.NewNatsConfig()
	s.Require().NoError(err)
	s.Equal("nats://nats:4222", nc.Url)

	cc, err := natsx.NewConsumerConfig()
	s.Require().NoError(err)
	s.Equal(natsx.ConsumerConfig{StreamName: "in-stream", DurableName: "in-consumer", MaxMessages: 1}, *cc)

	pc, err := natsx.NewPublisherConfig()
	s.Require().NoError(err)
	s.Equal(natsx.PublisherConfig{StreamName: "out-stream", Subject: "out.subject"}, *pc)
}

func (s *ConfigTestSuite) TestConsumerConfigRequiresChartKeys() {
	_, err := natsx.NewConsumerConfig()

	s.Error(err)
}
