package processor

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/quimnieto/akeneo-serverless-connector/aws/internal/models"
)

type EventProcessor interface {
	Parse(payload []byte) (*models.AkeneoEvent, error)
	Validate(event *models.AkeneoEvent) error
}

type processor struct{}

func NewProcessor() EventProcessor {
	return &processor{}
}

func (p *processor) Parse(payload []byte) (*models.AkeneoEvent, error) {
	if len(payload) == 0 {
		return nil, fmt.Errorf("payload is empty")
	}

	var event models.AkeneoEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return nil, fmt.Errorf("failed to parse JSON payload: %w", err)
	}

	return &event, nil
}

func (p *processor) Validate(event *models.AkeneoEvent) error {
	if event == nil {
		return fmt.Errorf("event is nil")
	}

	var missingFields []string

	if strings.TrimSpace(event.EventID) == "" {
		missingFields = append(missingFields, "event_id")
	}

	if strings.TrimSpace(event.EventType) == "" {
		missingFields = append(missingFields, "event_type")
	}

	if strings.TrimSpace(event.Timestamp) == "" {
		missingFields = append(missingFields, "timestamp")
	}

	if len(missingFields) > 0 {
		return fmt.Errorf("missing required fields: %s", strings.Join(missingFields, ", "))
	}

	return nil
}
