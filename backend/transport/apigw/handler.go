package apigw

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/aws/aws-lambda-go/events"
)

type Request = events.APIGatewayProxyRequest
type Response = events.APIGatewayProxyResponse

func NewResponse(body []byte) (Response, error) {
	return Response{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(body),
	}, nil
}

func NewErrorResponse(statusCode int, message string) (Response, error) {
	return Response{
		StatusCode: statusCode,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: `{"error":"` + message + `"}`,
	}, nil
}

func NewRequestHandler[Req any, Resp any](callback func(ctx context.Context, req Req) (Resp, error)) func(context.Context, Request) (Response, error) {
	return func(ctx context.Context, req Request) (Response, error) {
		slog.Info("Handling request: " + req.Body)

		var requestData Req
		if err := json.Unmarshal([]byte(req.Body), &requestData); err != nil {
			return NewErrorResponse(400, "Invalid request format: "+err.Error())
		}

		respData, err := callback(ctx, requestData)
		if err != nil {
			return NewErrorResponse(500, "Internal server error: "+err.Error())
		}

		respBytes, err := json.Marshal(respData)
		if err != nil {
			return NewErrorResponse(500, "Failed to marshal response: "+err.Error())
		}

		return NewResponse(respBytes)
	}
}
