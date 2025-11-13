package models

type AkeneoEvent struct {
	EventID   string                 `json:"event_id"`
	EventType string                 `json:"event_type"`
	Timestamp string                 `json:"timestamp"`
	Author    string                 `json:"author"`
	Data      map[string]interface{} `json:"data"`
}

type SNSMessage struct {
	Event      *AkeneoEvent      `json:"event"`
	ReceivedAt string            `json:"received_at"`
	Source     string            `json:"source"`
	Metadata   map[string]string `json:"metadata"`
}

