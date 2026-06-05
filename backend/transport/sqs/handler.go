package sqs

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"

	"github.com/aws/aws-lambda-go/events"
)

func NewRequestHandler[Req any, Resp any](callback func(ctx context.Context, req Req) (Resp, error)) func(context.Context, events.SQSEvent) error {
	return func(ctx context.Context, event events.SQSEvent) error {
		if len(event.Records) == 0 {
			return errors.New("no SQS record received")
		}
		if len(event.Records) > 1 {
			return errors.New("expected exactly one SQS record")
		}

		record := event.Records[0]
		var requestData Req
		if err := json.Unmarshal([]byte(record.Body), &requestData); err != nil {
			slog.Error("Failed to unmarshal SQS message", "error", err, "message_id", record.MessageId)
			return err
		}

		_, err := callback(ctx, requestData)
		if err != nil {
			slog.Error("Failed to process SQS message", "error", err, "message_id", record.MessageId)
			return err
		}

		return nil
	}
}
