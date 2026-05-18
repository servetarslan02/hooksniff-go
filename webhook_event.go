package hooksniff

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// WebhookEvent represents a parsed webhook event from HookSniff.
type WebhookEvent struct {
	// Event is the event type name (e.g., "endpoint.created")
	Event string `json:"event"`
	// Data is the event payload data
	Data map[string]interface{} `json:"data"`
	// Timestamp is the ISO 8601 timestamp string
	Timestamp string `json:"timestamp"`
}

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
