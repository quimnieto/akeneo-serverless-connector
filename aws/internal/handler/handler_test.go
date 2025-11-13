package handler_test

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"github.com/quimnieto/akeneo-serverless-connector/aws/internal/handler"
	"github.com/quimnieto/akeneo-serverless-connector/aws/internal/logger"
	"github.com/quimnieto/akeneo-serverless-connector/aws/internal/models"
	"github.com/quimnieto/akeneo-serverless-connector/aws/internal/processor"
)

type stubPublisher struct {
	lastEvent *models.AkeneoEvent
	err       error
}

func (s *stubPublisher) Publish(ctx context.Context, event *models.AkeneoEvent) error {
	if s.err != nil {
		return s.err
	}
	s.lastEvent = event
	return nil
}

func TestHandleSuccess(t *testing.T) {
	publisher := &stubPublisher{}
	h := handler.NewLambdaHandler(processor.NewProcessor(), publisher, logger.New())

	req := events.APIGatewayProxyRequest{
		Body: `{"event_id":"123","event_type":"product.updated","timestamp":"2024-10-01T10:00:00Z","author":"alice","data":{"sku":"ABC"}}`,
	}

	resp, err := h.Handle(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusAccepted, resp.StatusCode)

	var payload map[string]string
	assert.NoError(t, json.Unmarshal([]byte(resp.Body), &payload))
	assert.Equal(t, "accepted", payload["status"])
	assert.Equal(t, "123", payload["event_id"])

	assert.NotNil(t, publisher.lastEvent)
	assert.Equal(t, "product.updated", publisher.lastEvent.EventType)
}

func TestHandleBase64Payload(t *testing.T) {
	publisher := &stubPublisher{}
	h := handler.NewLambdaHandler(processor.NewProcessor(), publisher, logger.New())

	raw := `{"event_id":"abc","event_type":"product.created","timestamp":"2024-10-01T10:00:00Z","author":"bob","data":{"sku":"XYZ"}}`
	req := events.APIGatewayProxyRequest{
		Body:            base64.StdEncoding.EncodeToString([]byte(raw)),
		IsBase64Encoded: true,
	}

	resp, err := h.Handle(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusAccepted, resp.StatusCode)
	assert.NotNil(t, publisher.lastEvent)
	assert.Equal(t, "abc", publisher.lastEvent.EventID)
}

func TestHandleInvalidJSON(t *testing.T) {
	publisher := &stubPublisher{}
	h := handler.NewLambdaHandler(processor.NewProcessor(), publisher, logger.New())

	req := events.APIGatewayProxyRequest{
		Body: `{"event_id":}`,
	}

	resp, err := h.Handle(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var payload map[string]string
	assert.NoError(t, json.Unmarshal([]byte(resp.Body), &payload))
	assert.Equal(t, "INVALID_PAYLOAD", payload["error_code"])
}

func TestHandleValidationError(t *testing.T) {
	publisher := &stubPublisher{}
	h := handler.NewLambdaHandler(processor.NewProcessor(), publisher, logger.New())

	req := events.APIGatewayProxyRequest{
		Body: `{"event_id":"","event_type":"product.updated","timestamp":"2024-10-01T10:00:00Z","author":"alice","data":{"sku":"ABC"}}`,
	}

	resp, err := h.Handle(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var payload map[string]string
	assert.NoError(t, json.Unmarshal([]byte(resp.Body), &payload))
	assert.Equal(t, "MISSING_FIELD", payload["error_code"])
}

func TestHandlePublishError(t *testing.T) {
	publisher := &stubPublisher{
		err: errors.New("sns failure"),
	}
	h := handler.NewLambdaHandler(processor.NewProcessor(), publisher, logger.New())

	req := events.APIGatewayProxyRequest{
		Body: `{"event_id":"123","event_type":"product.updated","timestamp":"2024-10-01T10:00:00Z","author":"alice","data":{"sku":"ABC"}}`,
	}

	resp, err := h.Handle(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	var payload map[string]string
	assert.NoError(t, json.Unmarshal([]byte(resp.Body), &payload))
	assert.Equal(t, "SNS_PUBLISH_FAILED", payload["error_code"])
}

func TestHandleBase64DecodeError(t *testing.T) {
	publisher := &stubPublisher{}
	h := handler.NewLambdaHandler(processor.NewProcessor(), publisher, logger.New())

	req := events.APIGatewayProxyRequest{
		Body:            "not-base64",
		IsBase64Encoded: true,
	}

	resp, err := h.Handle(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var payload map[string]string
	assert.NoError(t, json.Unmarshal([]byte(resp.Body), &payload))
	assert.Equal(t, "INVALID_PAYLOAD", payload["error_code"])
}
