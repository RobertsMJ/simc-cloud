package natsx

import "github.com/RobertsMJ/simc-cloud/backend/config"

type NatsConfig struct {
	Url string `env:"QUEUE_URL,required"`
}

func NewNatsConfig() (*NatsConfig, error) {
	return config.Parse[NatsConfig]("")
}

type ConsumerConfig struct {
	StreamName  string `env:"STREAM_NAME,required"`
	DurableName string `env:"NAME,required"`
	MaxMessages int    `env:"MAX_MESSAGES" envDefault:"1"`
}

func NewConsumerConfig() (*ConsumerConfig, error) {
	return config.Parse[ConsumerConfig]("CONSUMER_")
}

type PublisherConfig struct {
	StreamName string `env:"STREAM_NAME,required"`
	Subject    string `env:"SUBJECT,required"`
}

func NewPublisherConfig() (*PublisherConfig, error) {
	return config.Parse[PublisherConfig]("PUBLISHER_")
}
