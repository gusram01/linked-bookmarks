package domain

import "encoding/json"

type EventAttributes struct {
	HTTPRequest HTTPRequest `json:"http_request"`
}

type HTTPRequest struct {
	ClientIP  string `json:"client_ip"`
	UserAgent string `json:"user_agent"`
}

type UserWebhookEvent struct {
	Data            json.RawMessage `json:"data"`
	EventAttributes EventAttributes `json:"event_attributes"`
	InstanceID      string          `json:"instance_id"`
	Object          string          `json:"object"`
	Timestamp       int64           `json:"timestamp"`
	Type            string          `json:"type"`
}
