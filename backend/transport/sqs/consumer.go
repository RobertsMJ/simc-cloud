package sqs

import (
	"context"
	"encoding/json"
	"log/slog"
	"sync"

	"github.com/aws/aws-lambda-go/events"
)

func NewConsumer[Message any](handler func(ctx context.Context, msg Message) error) func(context.Context, events.SQSEvent) (events.SQSEventResponse, error) {
	return func(ctx context.Context, event events.SQSEvent) (events.SQSEventResponse, error) {
		batchResponse := events.SQSEventResponse{}
		failedMessageIDs := make(chan string, len(event.Records))
		wg := sync.WaitGroup{}
		for _, record := range event.Records {
			wg.Go(func() {
				var msgData Message
				if err := json.Unmarshal([]byte(record.Body), &msgData); err != nil {
					slog.Error("Failed to unmarshal SQS message", "error", err, "message_id", record.MessageId)
					failedMessageIDs <- record.MessageId
					return
				}
				if err := handler(ctx, msgData); err != nil {
					slog.Error("Failed to process SQS message", "error", err, "message_id", record.MessageId)
					failedMessageIDs <- record.MessageId
					return
				}
			})
		}
		go func() {
			wg.Wait()
			close(failedMessageIDs)
		}()
		for failedMessageId := range failedMessageIDs {
			batchResponse.BatchItemFailures = append(batchResponse.BatchItemFailures, events.SQSBatchItemFailure{
				ItemIdentifier: failedMessageId,
			})
		}
		return batchResponse, nil
	}
}
