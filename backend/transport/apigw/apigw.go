package apigw

import "github.com/aws/aws-lambda-go/events"

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
