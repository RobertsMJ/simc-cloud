#!/bin/bash
set -euo pipefail

echo "Creating SQS queues..."
awslocal sqs create-queue --queue-name "$SIMULATION_QUEUE_NAME"
awslocal sqs create-queue --queue-name "$RESULTS_QUEUE_NAME"

echo "Creating DynamoDB table..."
awslocal dynamodb create-table \
  --table-name "$SIM_TABLE_NAME" \
  --attribute-definitions \
    AttributeName=PK,AttributeType=S \
    AttributeName=SK,AttributeType=S \
  --key-schema \
    AttributeName=PK,KeyType=HASH \
    AttributeName=SK,KeyType=RANGE \
  --billing-mode PAY_PER_REQUEST

echo "Local AWS resources ready."
