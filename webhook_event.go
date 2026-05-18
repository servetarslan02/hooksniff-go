package hooksniff

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// ─── Event Data Structs ─────────────────────────────────────────────

// EndpointCreatedEventData is the data payload for endpoint.created events.
type EndpointCreatedEventData struct {
	AppID      string  `json:"appId"`
	EndpointID string  `json:"endpointId"`
	AppUID     *string `json:"appUid,omitempty"`
}

// EndpointUpdatedEventData is the data payload for endpoint.updated events.
type EndpointUpdatedEventData struct {
	AppID      string  `json:"appId"`
	EndpointID string  `json:"endpointId"`
	AppUID     *string `json:"appUid,omitempty"`
}

// EndpointDeletedEventData is the data payload for endpoint.deleted events.
type EndpointDeletedEventData struct {
	AppID      string  `json:"appId"`
	EndpointID string  `json:"endpointId"`
	AppUID     *string `json:"appUid,omitempty"`
}

// EndpointEnabledEventData is the data payload for endpoint.enabled events.
type EndpointEnabledEventData struct {
	AppID      string  `json:"appId"`
	EndpointID string  `json:"endpointId"`
	AppUID     *string `json:"appUid,omitempty"`
}

// EndpointDisabledEventData is the data payload for endpoint.disabled events.
type EndpointDisabledEventData struct {
	AppID      string  `json:"appId"`
	EndpointID string  `json:"endpointId"`
	AppUID     *string `json:"appUid,omitempty"`
	FailSince  *string `json:"failSince,omitempty"`
	Trigger    *string `json:"trigger,omitempty"` // "none" | "first-failure" | "repeated-failure"
}

// LastAttemptInfo contains info about the last delivery attempt.
type LastAttemptInfo struct {
	ID                 string `json:"id"`
	Timestamp          string `json:"timestamp"`
	ResponseStatusCode int    `json:"responseStatusCode"`
}

// AttemptInfo contains info about a delivery attempt.
type AttemptInfo struct {
	ID                 string `json:"id"`
	Timestamp          string `json:"timestamp"`
	ResponseStatusCode int    `json:"responseStatusCode"`
}

// MessageAttemptExhaustedEventData is the data payload for message.attempt.exhausted events.
type MessageAttemptExhaustedEventData struct {
	AppID       string          `json:"appId"`
	MsgID       string          `json:"msgId"`
	LastAttempt LastAttemptInfo `json:"lastAttempt"`
	AppUID      *string         `json:"appUid,omitempty"`
}

// MessageAttemptFailingEventData is the data payload for message.attempt.failing events.
type MessageAttemptFailingEventData struct {
	AppID   string      `json:"appId"`
	MsgID   string      `json:"msgId"`
	Attempt AttemptInfo `json:"attempt"`
	AppUID  *string     `json:"appUid,omitempty"`
}

// MessageAttemptRecoveredEventData is the data payload for message.attempt.recovered events.
type MessageAttemptRecoveredEventData struct {
	AppID   string      `json:"appId"`
	MsgID   string      `json:"msgId"`
	Attempt AttemptInfo `json:"attempt"`
	AppUID  *string     `json:"appUid,omitempty"`
}

// ─── WebhookEvent (generic, backward compatible) ────────────────────

// WebhookEvent represents a parsed webhook event from HookSniff.
type WebhookEvent struct {
	// Event is the event type name (e.g., "endpoint.created")
	Event string `json:"event"`
	// Data is the event payload data (raw JSON — use TypedData() for typed access)
	Data map[string]interface{} `json:"data"`
	// Timestamp is the ISO 8601 timestamp string
	Timestamp string `json:"timestamp"`
}

// EventType returns the event type name.
func (e *WebhookEvent) EventType() string {
	return e.Event
}

// ─── Typed Data Parsing ─────────────────────────────────────────────

// ParseEndpointCreatedData parses the event data into EndpointCreatedEventData.
func (e *WebhookEvent) ParseEndpointCreatedData() (*EndpointCreatedEventData, error) {
	return parseEventData[EndpointCreatedEventData](e.Data)
}

// ParseEndpointUpdatedData parses the event data into EndpointUpdatedEventData.
func (e *WebhookEvent) ParseEndpointUpdatedData() (*EndpointUpdatedEventData, error) {
	return parseEventData[EndpointUpdatedEventData](e.Data)
}

// ParseEndpointDeletedData parses the event data into EndpointDeletedEventData.
func (e *WebhookEvent) ParseEndpointDeletedData() (*EndpointDeletedEventData, error) {
	return parseEventData[EndpointDeletedEventData](e.Data)
}

// ParseEndpointEnabledData parses the event data into EndpointEnabledEventData.
func (e *WebhookEvent) ParseEndpointEnabledData() (*EndpointEnabledEventData, error) {
	return parseEventData[EndpointEnabledEventData](e.Data)
}

// ParseEndpointDisabledData parses the event data into EndpointDisabledEventData.
func (e *WebhookEvent) ParseEndpointDisabledData() (*EndpointDisabledEventData, error) {
	return parseEventData[EndpointDisabledEventData](e.Data)
}

// ParseMessageAttemptExhaustedData parses the event data into MessageAttemptExhaustedEventData.
func (e *WebhookEvent) ParseMessageAttemptExhaustedData() (*MessageAttemptExhaustedEventData, error) {
	return parseEventData[MessageAttemptExhaustedEventData](e.Data)
}

// ParseMessageAttemptFailingData parses the event data into MessageAttemptFailingEventData.
func (e *WebhookEvent) ParseMessageAttemptFailingData() (*MessageAttemptFailingEventData, error) {
	return parseEventData[MessageAttemptFailingEventData](e.Data)
}

// ParseMessageAttemptRecoveredData parses the event data into MessageAttemptRecoveredEventData.
func (e *WebhookEvent) ParseMessageAttemptRecoveredData() (*MessageAttemptRecoveredEventData, error) {
	return parseEventData[MessageAttemptRecoveredEventData](e.Data)
}

// parseEventData is a generic helper that marshals the raw map into JSON, then unmarshals into T.
func parseEventData[T any](data map[string]interface{}) (*T, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal event data: %w", err)
	}
	var result T
	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		return nil, fmt.Errorf("failed to parse event data: %w", err)
	}
	return &result, nil
}

// ─── Verify + Parse ─────────────────────────────────────────────────

// VerifyAndParse validates the webhook signature and parses the payload into a WebhookEvent.
//
// Returns the parsed WebhookEvent on success, or an error if the signature
// is invalid, the timestamp is outside tolerance, or the payload cannot be parsed.
func (wh *Webhook) VerifyAndParse(payload []byte, headers http.Header) (*WebhookEvent, error) {
	if err := wh.verify(payload, headers, true); err != nil {
		return nil, err
	}

	var event WebhookEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return nil, fmt.Errorf("failed to parse webhook payload: %w", err)
	}

	return &event, nil
}

// VerifyAndParseIgnoringTimestamp validates the webhook signature (without timestamp check)
// and parses the payload into a WebhookEvent.
//
// WARNING: This function does not check the signature's timestamp.
// We recommend using VerifyAndParse instead.
func (wh *Webhook) VerifyAndParseIgnoringTimestamp(payload []byte, headers http.Header) (*WebhookEvent, error) {
	if err := wh.verify(payload, headers, false); err != nil {
		return nil, err
	}

	var event WebhookEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return nil, fmt.Errorf("failed to parse webhook payload: %w", err)
	}

	return &event, nil
}
