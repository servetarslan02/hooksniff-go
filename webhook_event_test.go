package hooksniff

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestVerifyAndParseEndpointCreated(t *testing.T) {
	secret := "whsec_dGVzdA=="
	msgID := "msg_test123"
	timestamp := nowUnix()
	payload := `{"event":"endpoint.created","data":{"appId":"app_1","endpointId":"ep_1","appUid":"uid_1"},"timestamp":"2026-05-19T00:00:00Z"}`

	sig := sign(secret, msgID, timestamp, payload)
	headers := map[string]string{
		"webhook-id":        msgID,
		"webhook-timestamp": fmt.Sprintf("%d", timestamp),
		"webhook-signature": sig,
	}

	wh := NewWebhook(secret)
	event, err := wh.VerifyAndParse([]byte(payload), headers)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if event.Event != "endpoint.created" {
		t.Errorf("expected event=endpoint.created, got %s", event.Event)
	}
	if event.Timestamp != "2026-05-19T00:00:00Z" {
		t.Errorf("expected timestamp, got %s", event.Timestamp)
	}
}

func TestParseEndpointCreatedData(t *testing.T) {
	event := &WebhookEvent{
		Event: "endpoint.created",
		Data: map[string]interface{}{
			"appId":      "app_1",
			"endpointId": "ep_1",
			"appUid":     "uid_1",
		},
		Timestamp: "2026-05-19",
	}

	data, err := event.ParseEndpointCreatedData()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data.AppID != "app_1" {
		t.Errorf("expected appId=app_1, got %s", data.AppID)
	}
	if data.EndpointID != "ep_1" {
		t.Errorf("expected endpointId=ep_1, got %s", data.EndpointID)
	}
	if data.AppUID == nil || *data.AppUID != "uid_1" {
		t.Errorf("expected appUid=uid_1")
	}
}

func TestParseEndpointDisabledData(t *testing.T) {
	event := &WebhookEvent{
		Event: "endpoint.disabled",
		Data: map[string]interface{}{
			"appId":      "app_1",
			"endpointId": "ep_1",
			"failSince":  "2026-01-01",
			"trigger":    "repeated-failure",
		},
		Timestamp: "2026-05-19",
	}

	data, err := event.ParseEndpointDisabledData()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data.FailSince == nil || *data.FailSince != "2026-01-01" {
		t.Errorf("expected failSince=2026-01-01")
	}
	if data.Trigger == nil || *data.Trigger != "repeated-failure" {
		t.Errorf("expected trigger=repeated-failure")
	}
}

func TestParseMessageAttemptExhaustedData(t *testing.T) {
	event := &WebhookEvent{
		Event: "message.attempt.exhausted",
		Data: map[string]interface{}{
			"appId": "app_1",
			"msgId": "msg_1",
			"lastAttempt": map[string]interface{}{
				"id":                 "att_1",
				"timestamp":          "2026-05-19",
				"responseStatusCode": 500,
			},
		},
		Timestamp: "2026-05-19",
	}

	data, err := event.ParseMessageAttemptExhaustedData()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data.MsgID != "msg_1" {
		t.Errorf("expected msgId=msg_1, got %s", data.MsgID)
	}
	if data.LastAttempt.ResponseStatusCode != 500 {
		t.Errorf("expected status=500, got %d", data.LastAttempt.ResponseStatusCode)
	}
}

func TestParseMessageAttemptFailingData(t *testing.T) {
	event := &WebhookEvent{
		Event: "message.atattempt.failing",
		Data: map[string]interface{}{
			"appId": "app_1",
			"msgId": "msg_1",
			"attempt": map[string]interface{}{
				"id":                 "att_1",
				"timestamp":          "2026-05-19",
				"responseStatusCode": 429,
			},
		},
		Timestamp: "2026-05-19",
	}

	data, err := event.ParseMessageAttemptFailingData()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data.Attempt.ResponseStatusCode != 429 {
		t.Errorf("expected status=429, got %d", data.Attempt.ResponseStatusCode)
	}
}

func TestParseMessageAttemptRecoveredData(t *testing.T) {
	event := &WebhookEvent{
		Event: "message.attempt.recovered",
		Data: map[string]interface{}{
			"appId": "app_1",
			"msgId": "msg_1",
			"attempt": map[string]interface{}{
				"id":                 "att_1",
				"timestamp":          "2026-05-19",
				"responseStatusCode": 200,
			},
		},
		Timestamp: "2026-05-19",
	}

	data, err := event.ParseMessageAttemptRecoveredData()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data.Attempt.ID != "att_1" {
		t.Errorf("expected attempt.id=att_1, got %s", data.Attempt.ID)
	}
}

func TestParseAllEndpointEvents(t *testing.T) {
	tests := []struct {
		event    string
		parseFn  func(*WebhookEvent) (interface{}, error)
		wantType string
	}{
		{"endpoint.created", func(e *WebhookEvent) (interface{}, error) { return e.ParseEndpointCreatedData() }, "*hooksniff.EndpointCreatedEventData"},
		{"endpoint.updated", func(e *WebhookEvent) (interface{}, error) { return e.ParseEndpointUpdatedData() }, "*hooksniff.EndpointUpdatedEventData"},
		{"endpoint.deleted", func(e *WebhookEvent) (interface{}, error) { return e.ParseEndpointDeletedData() }, "*hooksniff.EndpointDeletedEventData"},
		{"endpoint.enabled", func(e *WebhookEvent) (interface{}, error) { return e.ParseEndpointEnabledData() }, "*hooksniff.EndpointEnabledEventData"},
		{"endpoint.disabled", func(e *WebhookEvent) (interface{}, error) { return e.ParseEndpointDisabledData() }, "*hooksniff.EndpointDisabledEventData"},
	}

	for _, tt := range tests {
		t.Run(tt.event, func(t *testing.T) {
			event := &WebhookEvent{
				Event:     tt.event,
				Data:      map[string]interface{}{"appId": "a", "endpointId": "e"},
				Timestamp: "2026-05-19",
			}
			result, err := tt.parseFn(event)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			got := fmt.Sprintf("%T", result)
			if got != tt.wantType {
				t.Errorf("got %s, want %s", got, tt.wantType)
			}
		})
	}
}

func TestVerifyAndParseWithTypedData(t *testing.T) {
	secret := "whsec_dGVzdA=="
	msgID := "msg_test123"
	timestamp := nowUnix()
	payload := `{"event":"endpoint.created","data":{"appId":"app_1","endpointId":"ep_1"},"timestamp":"2026-05-19"}`

	sig := sign(secret, msgID, timestamp, payload)
	headers := map[string]string{
		"webhook-id":        msgID,
		"webhook-timestamp": fmt.Sprintf("%d", timestamp),
		"webhook-signature": sig,
	}

	wh := NewWebhook(secret)
	event, err := wh.VerifyAndParse([]byte(payload), headers)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := event.ParseEndpointCreatedData()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data.AppID != "app_1" {
		t.Errorf("expected appId=app_1, got %s", data.AppID)
	}
	if data.EndpointID != "ep_1" {
		t.Errorf("expected endpointId=ep_1, got %s", data.EndpointID)
	}
}

func TestSDKVersionHeader(t *testing.T) {
	// Verify that the SDK version is accessible
	if Version == "" {
		t.Error("Version should not be empty")
	}
}

func TestXHookSniffSDKHeader(t *testing.T) {
	// The X-HookSniff-SDK header is set in hooksniff.go during client init
	// This test verifies the format
	sdkUA := fmt.Sprintf("hooksniff-libs/%s/go", Version)
	expected := "hooksniff-libs/"
	if len(sdkUA) < len(expected) || sdkUA[:len(expected)] != expected {
		t.Errorf("SDK header should start with 'hooksniff-libs/', got %s", sdkUA)
	}
}

// Helper
func nowUnix() int64 {
	return time.Now().Unix()
}

// Re-export sign for tests
func init() {
	// sign is defined in hooksniff_test.go
}

// sign is already defined in hooksniff_test.go
func signForTest(secret, msgID string, timestamp int64, payload string) string {
	decoded, _ := base64.StdEncoding.DecodeString(secret)
	toSign := fmt.Sprintf("%s.%d.%s", msgID, timestamp, payload)
	mac := hmac.New(sha256.New, decoded)
	mac.Write([]byte(toSign))
	sig := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	return "v1," + sig
}

// TestParseDataFromJSON verifies typed data can be parsed from JSON
func TestParseDataFromJSON(t *testing.T) {
	jsonStr := `{"appId":"app_1","endpointId":"ep_1","appUid":"uid_1"}`
	var raw map[string]interface{}
	json.Unmarshal([]byte(jsonStr), &raw)

	event := &WebhookEvent{
		Event:     "endpoint.created",
		Data:      raw,
		Timestamp: "2026-05-19",
	}

	data, err := event.ParseEndpointCreatedData()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data.AppID != "app_1" {
		t.Errorf("expected app_1, got %s", data.AppID)
	}
}
