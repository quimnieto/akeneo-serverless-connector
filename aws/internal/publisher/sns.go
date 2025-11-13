package publisher

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sns/types"
	"github.com/quimnieto/akeneo-serverless-connector/aws/internal/models"
)

type SNSPublisher interface {
	Publish(ctx context.Context, event *models.AkeneoEvent) error
}

type SNSClient interface {
	Publish(ctx context.Context, params *sns.PublishInput, optFns ...func(*sns.Options)) (*sns.PublishOutput, error)
}

type snsPublisher struct {
	client   SNSClient
	topicARN string
}

func NewSNSPublisher(client *sns.Client, topicARN string) SNSPublisher {
	return &snsPublisher{
		client:   client,
		topicARN: topicARN,
	}
}

func NewSNSPublisherWithClient(client SNSClient, topicARN string) SNSPublisher {
	return &snsPublisher{
		client:   client,
		topicARN: topicARN,
	}
}

func (p *snsPublisher) Publish(ctx context.Context, event *models.AkeneoEvent) error {
	if event == nil {
		return fmt.Errorf("event is nil")
	}

	snsMessage := &models.SNSMessage{
		Event:      event,
		ReceivedAt: time.Now().UTC().Format(time.RFC3339),
		Source:     "akeneo-webhook",
		Metadata: map[string]string{
			"event_id":   event.EventID,
			"event_type": event.EventType,
		},
	}

	messageBody, err := json.Marshal(snsMessage)
	if err != nil {
		return fmt.Errorf("failed to marshal SNS message: %w", err)
	}

	messageAttributes := map[string]types.MessageAttributeValue{
		"event_type": {
			DataType:    aws.String("String"),
			StringValue: aws.String(event.EventType),
		},
		"timestamp": {
			DataType:    aws.String("String"),
			StringValue: aws.String(event.Timestamp),
		},
		"event_id": {
			DataType:    aws.String("String"),
			StringValue: aws.String(event.EventID),
		},
	}

	maxRetries := 3
	backoffDurations := []time.Duration{
		1 * time.Second,
		2 * time.Second,
		4 * time.Second,
	}

	var lastErr error
	for attempt := 0; attempt < maxRetries; attempt++ {
		_, err := p.client.Publish(ctx, &sns.PublishInput{
			TopicArn:          aws.String(p.topicARN),
			Message:           aws.String(string(messageBody)),
			MessageAttributes: messageAttributes,
		})

		if err == nil {
			return nil
		}

		lastErr = err

		if attempt < maxRetries-1 {
			select {
			case <-ctx.Done():
				return fmt.Errorf("context cancelled during retry: %w", ctx.Err())
			case <-time.After(backoffDurations[attempt]):
			}
		}
	}

	return fmt.Errorf("failed to publish to SNS after %d attempts: %w", maxRetries, lastErr)
}

