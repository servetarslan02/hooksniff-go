package models

type PollingEndpointMessageOut struct {
	EventId   string `json:"eventId"`
	EventType string `json:"eventType"`
	Timestamp string `json:"timestamp"`
}
