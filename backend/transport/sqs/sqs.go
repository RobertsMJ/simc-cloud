package sqs

import (
	"context"
	"reflect"

	"github.com/RobertsMJ/simc-cloud-backend/logger"
	"github.com/RobertsMJ/simc-cloud-backend/simc"
	"github.com/aws/aws-lambda-go/events"
)

type Request events.SQSEvent

// SQS handlers just return an error
type Response = error

func NewRequestHandler[Req simc.Unmarshaler, Resp simc.Marshaler](callback func(ctx context.Context, req Req) (Resp, error)) func(context.Context, Request) (Response, error) {
	return func(ctx context.Context, event Request) (Response, error) {
		for _, record := range event.Records {
			var requestData Req
			requestData = reflect.New(reflect.TypeOf(requestData).Elem()).Interface().(Req)
			if err := requestData.UnmarshalSimC([]byte(record.Body)); err != nil {
				logger.Error("Failed to unmarshal SQS message", "error", err, "message_id", record.MessageId)
				continue
			}

			// Call the callback with the parsed request
			_, err := callback(ctx, requestData)
			if err != nil {
				logger.Error("Failed to process SQS message", "error", err, "message_id", record.MessageId)
				continue
			}
		}

		return nil, nil
	}
}
