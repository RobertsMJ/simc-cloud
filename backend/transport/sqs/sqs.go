package sqs

import "github.com/aws/aws-lambda-go/events"

type Request events.SQSEvent

// SQS handlers just return an error
type Response = error
