package sqs

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/RobertsMJ/simc-cloud-backend/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type Publisher[T any] struct {
	client   *sqs.Client
	queueURL string
}

func NewPublisher[Message any](ctx context.Context, queueName string) *Publisher[Message] {
	client := sqs.NewFromConfig(config.LoadAWS(ctx))
	queueURL, err := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
		QueueName: aws.String(queueName),
	})
	if err != nil {
		panic(fmt.Errorf("failed to get queue URL: %w", err))
	}
	return &Publisher[Message]{client: client, queueURL: *queueURL.QueueUrl}
}

func (p *Publisher[Message]) Publish(ctx context.Context, msg Message) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	_, err = p.client.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    aws.String(p.queueURL),
		MessageBody: aws.String(string(body)),
	})
	if err != nil {
		return fmt.Errorf("failed to publish message to SQS: %w", err)
	}

	return nil
}
