package sqs

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/RobertsMJ/simc-cloud-backend/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type Publisher[T any] struct {
	client   *sqs.Client
	queueURL string
}

func NewPublisher[Message any](ctx context.Context, queueURL string) *Publisher[Message] {
	client := NewClient(ctx, config.LoadAWS(ctx))
	return &Publisher[Message]{client: client, queueURL: queueURL}
}

func (p *Publisher[Message]) Publish(ctx context.Context, msg Message) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	slog.Debug("Publishing message to SQS", "queueURL", p.queueURL, "message", string(body))
	_, err = p.client.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    aws.String(p.queueURL),
		MessageBody: aws.String(string(body)),
	})
	if err != nil {
		return fmt.Errorf("failed to publish message to SQS: %w", err)
	}

	return nil
}
