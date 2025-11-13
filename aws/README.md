# Akeneo Serverless Connector - AWS Lambda

AWS Lambda implementation of the Akeneo webhook connector using Go.

## Overview

This serverless function receives webhook events from Akeneo PIM via API Gateway, validates the payload, and publishes events to an SNS topic for downstream processing.

## Architecture

```
Akeneo PIM → API Gateway → Lambda Function → SNS Topic → Downstream Services
```

The Lambda function:
- Receives webhook events from Akeneo PIM via API Gateway
- Validates the event payload structure
- Publishes validated events to an SNS topic with retry logic
- Returns appropriate HTTP responses

## Features

- **Serverless Architecture**: Built for AWS Lambda with minimal cold start time
- **Event Validation**: Ensures all required fields are present before publishing
- **SNS Integration**: Publishes events to SNS with message attributes for filtering
- **Retry Logic**: Automatic retry with exponential backoff for SNS publishing
- **Structured Logging**: JSON-formatted logs for CloudWatch integration
- **Error Handling**: Comprehensive error handling with detailed error codes
- **Base64 Support**: Handles both plain and base64-encoded payloads

## Project Structure

```
.
├── cmd/
│   └── lambda/
│       └── main.go              # Lambda entry point
├── internal/
│   ├── errors/
│   │   └── errors.go            # Custom error types
│   ├── handler/
│   │   ├── handler.go           # Main Lambda handler
│   │   └── handler_test.go      # Handler tests
│   ├── logger/
│   │   └── logger.go            # Structured logging
│   ├── models/
│   │   └── event.go             # Event data models
│   ├── processor/
│   │   └── processor.go         # Event parsing and validation
│   └── publisher/
│       └── sns.go               # SNS publishing logic
├── go.mod
└── go.sum
```

## Prerequisites

- Go 1.23 or later
- AWS CLI configured with appropriate credentials
- AWS account with permissions for Lambda, SNS, and API Gateway

## Installation

1. Navigate to the AWS directory:
```bash
cd aws
```

2. Install Go dependencies:
```bash
go mod download
```

## Building

Build the Lambda function:

```bash
GOOS=linux GOARCH=amd64 go build -o bootstrap cmd/lambda/main.go
zip function.zip bootstrap
```

## Configuration

The Lambda function requires the following environment variable:

- `SNS_TOPIC_ARN`: The ARN of the SNS topic where events will be published
- `LOG_LEVEL` (optional): Logging level (DEBUG, INFO, ERROR). Default: INFO

## Deployment

### Manual Deployment

1. Create an SNS topic:
```bash
aws sns create-topic --name akeneo-events
```

2. Create the Lambda function:
```bash
aws lambda create-function \
  --function-name akeneo-webhook-handler \
  --runtime provided.al2 \
  --handler bootstrap \
  --zip-file fileb://function.zip \
  --role arn:aws:iam::YOUR_ACCOUNT:role/lambda-execution-role \
  --environment Variables={SNS_TOPIC_ARN=arn:aws:sns:REGION:ACCOUNT:akeneo-events}
```

3. Create an API Gateway REST API and integrate it with the Lambda function

### IAM Permissions

The Lambda execution role needs the following permissions:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "sns:Publish"
      ],
      "Resource": "arn:aws:sns:REGION:ACCOUNT:akeneo-events"
    },
    {
      "Effect": "Allow",
      "Action": [
        "logs:CreateLogGroup",
        "logs:CreateLogStream",
        "logs:PutLogEvents"
      ],
      "Resource": "arn:aws:logs:*:*:*"
    }
  ]
}
```

## Event Format

### Input (Akeneo Webhook)

```json
{
  "event_id": "unique-event-id",
  "event_type": "product.updated",
  "timestamp": "2024-10-01T10:00:00Z",
  "author": "username",
  "data": {
    "sku": "PRODUCT-SKU",
    "additional": "fields"
  }
}
```

### Output (SNS Message)

```json
{
  "event": {
    "event_id": "unique-event-id",
    "event_type": "product.updated",
    "timestamp": "2024-10-01T10:00:00Z",
    "author": "username",
    "data": {
      "sku": "PRODUCT-SKU"
    }
  },
  "received_at": "2024-10-01T10:00:01Z",
  "source": "akeneo-webhook",
  "metadata": {
    "event_id": "unique-event-id",
    "event_type": "product.updated"
  }
}
```

## API Responses

### Success (202 Accepted)

```json
{
  "status": "accepted",
  "event_id": "unique-event-id"
}
```

### Error (4xx/5xx)

```json
{
  "error_code": "INVALID_PAYLOAD",
  "message": "Invalid webhook payload",
  "details": "specific error details"
}
```

### Error Codes

- `INVALID_PAYLOAD`: Malformed JSON or invalid payload structure
- `MISSING_FIELD`: Required field is missing from the payload
- `SNS_PUBLISH_FAILED`: Failed to publish event to SNS

## Testing

Run the test suite:

```bash
go test ./...
```

Run tests with coverage:

```bash
go test -cover ./...
```

Run tests with verbose output:

```bash
go test -v ./...
```

## Development

### Local Testing

You can test the handler locally by creating a test file:

```go
package main

import (
    "context"
    "fmt"
    "github.com/aws/aws-lambda-go/events"
    "github.com/quimnieto/akeneo-serverless-connector/aws/internal/handler"
    // ... other imports
)

func main() {
    // Initialize handler with mock dependencies
    h := handler.NewLambdaHandler(...)
    
    req := events.APIGatewayProxyRequest{
        Body: `{"event_id":"test","event_type":"product.updated","timestamp":"2024-10-01T10:00:00Z"}`,
    }
    
    resp, err := h.Handle(context.Background(), req)
    fmt.Printf("Response: %+v\nError: %v\n", resp, err)
}
```

## Logging

The function uses structured JSON logging compatible with CloudWatch Logs Insights:

```json
{
  "timestamp": "2024-10-01T10:00:00Z",
  "level": "INFO",
  "message": "event published",
  "request_id": "lambda-request-id",
  "fields": {
    "event_id": "unique-event-id",
    "event_type": "product.updated"
  }
}
```

### CloudWatch Logs Insights Queries

Find all errors:
```
fields @timestamp, message, error
| filter level = "ERROR"
| sort @timestamp desc
```

Track specific event:
```
fields @timestamp, message, fields.event_id
| filter fields.event_id = "your-event-id"
| sort @timestamp desc
```

## Monitoring

### CloudWatch Metrics

Key metrics to monitor:
- Lambda invocations
- Lambda errors
- Lambda duration
- SNS publish success/failure

### Alarms

Consider setting up CloudWatch alarms for:
- Error rate > threshold
- Duration > timeout threshold
- Throttles > 0

## Performance

- **Cold Start**: ~100-200ms (Go runtime)
- **Warm Execution**: ~10-50ms
- **Memory**: 128MB recommended (adjust based on payload size)
- **Timeout**: 30 seconds recommended

## Troubleshooting

### Common Issues

**Lambda timeout**
- Increase timeout setting
- Check SNS publish latency
- Review retry logic

**SNS publish failures**
- Verify IAM permissions
- Check SNS topic ARN
- Review CloudWatch logs

**Validation errors**
- Verify Akeneo webhook payload format
- Check required fields are present
- Review error details in response

## Cost Optimization

- Use ARM64 architecture for lower costs (update build command)
- Set appropriate memory allocation
- Enable Lambda reserved concurrency if needed
- Use SNS message filtering to reduce downstream processing

## Back to Main Documentation

See the [main README](../README.md) for multi-cloud overview and other implementations.

## License

This project is licensed under the MIT License - see the [LICENSE](../LICENSE) file for details.
