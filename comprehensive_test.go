package hooksniff

import (
	"encoding/json"
	"fmt"
	"testing"
)

// ═══════════════════════════════════════════════════════════════════
// WEBHOOK SIGNATURE TESTS
// ═══════════════════════════════════════════════════════════════════

func TestSignDeterministic(t *testing.T) {
	secret := "whsec_dGVzdA=="
	wh := NewWebhook(secret)
	ts := int64(1700000000)

	sig1, _ := wh.Sign("msg_1", ts, []byte("payload"))
	sig2, _ := wh.Sign("msg_1", ts, []byte("payload"))
	if sig1 != sig2 {
		t.Errorf("sign should be deterministic: %s != %s", sig1, sig2)
	}
}

func TestSignDifferentPayloads(t *testing.T) {
	secret := "whsec_dGVzdA=="
	wh := NewWebhook(secret)
	ts := int64(1700000000)

	sig1, _ := wh.Sign("msg_1", ts, []byte("payload1"))
	sig2, _ := wh.Sign("msg_1", ts, []byte("payload2"))
	if sig1 == sig2 {
		t.Error("different payloads should produce different signatures")
	}
}

func TestSignDifferentSecrets(t *testing.T) {
	wh1, _ := NewWebhook("whsec_dGVzdA==")
	wh2, _ := NewWebhook("whsec_b3RoZXI=")
	ts := int64(1700000000)

	sig1, _ := wh1.Sign("msg_1", ts, []byte("payload"))
	sig2, _ := wh2.Sign("msg_1", ts, []byte("payload"))
	if sig1 == sig2 {
		t.Error("different secrets should produce different signatures")
	}
}

// ═══════════════════════════════════════════════════════════════════
// VERIFY AND PARSE TESTS
// ═══════════════════════════════════════════════════════════════════

func TestVerifyAndParseValid(t *testing.T) {
	secret := "whsec_dGVzdA=="
	payload := `{"event":"endpoint.created","data":{"appId":"a1","endpointId":"e1"},"timestamp":"2026-05-19"}`
	ts := nowUnix()
	sig := sign(secret, "msg_1", ts, payload)

	wh := NewWebhook(secret)
	event, err := wh.VerifyAndParse([]byte(payload), map[string]string{
		"webhook-id":        "msg_1",
		"webhook-timestamp": fmt.Sprintf("%d", ts),
		"webhook-signature": sig,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if event.Event != "endpoint.created" {
		t.Errorf("expected endpoint.created, got %s", event.Event)
	}
	if event.Timestamp != "2026-05-19" {
		t.Errorf("expected 2026-05-19, got %s", event.Timestamp)
	}
}

func TestVerifyAndParseInvalidSignature(t *testing.T) {
	wh := NewWebhook("whsec_dGVzdA==")
	_, err := wh.VerifyAndParse([]byte(`{"event":"test"}`), map[string]string{
		"webhook-id":        "msg_1",
		"webhook-timestamp": fmt.Sprintf("%d", nowUnix()),
		"webhook-signature": "v1,invalid",
	})
	if err == nil {
		t.Error("expected error for invalid signature")
	}
}

func TestVerifyAndParseOldTimestamp(t *testing.T) {
	secret := "whsec_dGVzdA=="
	payload := `{"event":"test"}`
	ts := nowUnix() - 600
	sig := sign(secret, "msg_1", ts, payload)

	wh := NewWebhook(secret)
	_, err := wh.VerifyAndParse([]byte(payload), map[string]string{
		"webhook-id":        "msg_1",
		"webhook-timestamp": fmt.Sprintf("%d", ts),
		"webhook-signature": sig,
	})
	if err == nil {
		t.Error("expected error for old timestamp")
	}
}

func TestVerifyAndParseIgnoringTimestamp(t *testing.T) {
	secret := "whsec_dGVzdA=="
	payload := `{"event":"test"}`
	ts := nowUnix() - 600
	sig := sign(secret, "msg_1", ts, payload)

	wh := NewWebhook(secret)
	event, err := wh.VerifyAndParseIgnoringTimestamp([]byte(payload), map[string]string{
		"webhook-id":        "msg_1",
		"webhook-timestamp": fmt.Sprintf("%d", ts),
		"webhook-signature": sig,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if event.Event != "test" {
		t.Errorf("expected test, got %s", event.Event)
	}
}

func TestVerifyAndParseMissingHeaders(t *testing.T) {
	wh := NewWebhook("whsec_dGVzdA==")
	_, err := wh.VerifyAndParse([]byte(`{}`), map[string]string{})
	if err == nil {
		t.Error("expected error for missing headers")
	}
}

func TestVerifyAndParseHooksniffHeaders(t *testing.T) {
	secret := "whsec_dGVzdA=="
	payload := `{"event":"test"}`
	ts := nowUnix()
	sig := sign(secret, "msg_1", ts, payload)

	wh := NewWebhook(secret)
	event, err := wh.VerifyAndParse([]byte(payload), map[string]string{
		"hooksniff-id":        "msg_1",
		"hooksniff-timestamp": fmt.Sprintf("%d", ts),
		"hooksniff-signature": sig,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if event.Event != "test" {
		t.Errorf("expected test, got %s", event.Event)
	}
}

// ═══════════════════════════════════════════════════════════════════
// TYPED DATA PARSING TESTS
// ═══════════════════════════════════════════════════════════════════

func TestParseEndpointCreatedData(t *testing.T) {
	event := &WebhookEvent{
		Event:     "endpoint.created",
		Data:      map[string]interface{}{"appId": "a1", "endpointId": "e1", "appUid": "u1"},
		Timestamp: "2026-05-19",
	}
	data, err := event.ParseEndpointCreatedData()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data.AppID != "a1" {
		t.Errorf("expected a1, got %s", data.AppID)
	}
	if data.EndpointID != "e1" {
		t.Errorf("expected e1, got %s", data.EndpointID)
	}
	if data.AppUID == nil || *data.AppUID != "u1" {
		t.Errorf("expected u1")
	}
}

func TestParseEndpointUpdatedData(t *testing.T) {
	event := &WebhookEvent{
		Event:     "endpoint.updated",
		Data:      map[string]interface{}{"appId": "a1", "endpointId": "e1"},
		Timestamp: "",
	}
	data, err := event.ParseEndpointUpdatedData()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data.AppID != "a1" {
		t.Errorf("expected a1, got %s", data.AppID)
	}
}

func TestParseEndpointDeletedData(t *testing.T) {
	event := &WebhookEvent{
		Event:     "endpoint.deleted",
		Data:      map[string]interface{}{"appId": "a1", "endpointId": "e1"},
		Timestamp: "",
	}
	data, err := event.ParseEndpointDeletedData()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data.EndpointID != "e1" {
		t.Errorf("expected e1, got %s", data.EndpointID)
	}
}

func TestParseEndpointEnabledData(t *testing.T) {
	event := &WebhookEvent{
		Event:     "endpoint.enabled",
		Data:      map[string]interface{}{"appId": "a1", "endpointId": "e1"},
		Timestamp: "",
	}
	data, err := event.ParseEndpointEnabledData()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data.AppID != "a1" {
		t.Errorf("expected a1, got %s", data.AppID)
	}
}

func TestParseEndpointDisabledDataWithExtras(t *testing.T) {
	event := &WebhookEvent{
		Event: "endpoint.disabled",
		Data: map[string]interface{}{
			"appId": "a1", "endpointId": "e1",
			"failSince": "2026-01", "trigger": "repeated-failure",
		},
		Timestamp: "",
	}
	data, err := event.ParseEndpointDisabledData()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data.FailSince == nil || *data.FailSince != "2026-01" {
		t.Error("expected failSince=2026-01")
	}
	if data.Trigger == nil || *data.Trigger != "repeated-failure" {
		t.Error("expected trigger=repeated-failure")
	}
}

func TestParseEndpointDisabledDataNilOptionals(t *testing.T) {
	event := &WebhookEvent{
		Event:     "endpoint.disabled",
		Data:      map[string]interface{}{"appId": "a1", "endpointId": "e1"},
		Timestamp: "",
	}
	data, err := event.ParseEndpointDisabledData()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data.FailSince != nil {
		t.Error("expected nil failSince")
	}
	if data.Trigger != nil {
		t.Error("expected nil trigger")
	}
}

func TestParseMessageAttemptExhaustedData(t *testing.T) {
	event := &WebhookEvent{
		Event: "message.attempt.exhausted",
		Data: map[string]interface{}{
			"appId": "a1", "msgId": "m1",
			"lastAttempt": map[string]interface{}{
				"id": "att_1", "timestamp": "t", "responseStatusCode": 500,
			},
		},
		Timestamp: "",
	}
	data, err := event.ParseMessageAttemptExhaustedData()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data.MsgID != "m1" {
		t.Errorf("expected m1, got %s", data.MsgID)
	}
	if data.LastAttempt.ID != "att_1" {
		t.Errorf("expected att_1, got %s", data.LastAttempt.ID)
	}
	if data.LastAttempt.ResponseStatusCode != 500 {
		t.Errorf("expected 500, got %d", data.LastAttempt.ResponseStatusCode)
	}
}

func TestParseMessageAttemptFailingData(t *testing.T) {
	event := &WebhookEvent{
		Event: "message.atattempt.failing",
		Data: map[string]interface{}{
			"appId": "a1", "msgId": "m1",
			"attempt": map[string]interface{}{
				"id": "att_1", "timestamp": "t", "responseStatusCode": 429,
			},
		},
		Timestamp: "",
	}
	data, err := event.ParseMessageAttemptFailingData()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data.Attempt.ResponseStatusCode != 429 {
		t.Errorf("expected 429, got %d", data.Attempt.ResponseStatusCode)
	}
}

func TestParseMessageAttemptRecoveredData(t *testing.T) {
	event := &WebhookEvent{
		Event: "message.attempt.recovered",
		Data: map[string]interface{}{
			"appId": "a1", "msgId": "m1",
			"attempt": map[string]interface{}{
				"id": "att_1", "timestamp": "t", "responseStatusCode": 200,
			},
		},
		Timestamp: "",
	}
	data, err := event.ParseMessageAttemptRecoveredData()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data.Attempt.ResponseStatusCode != 200 {
		t.Errorf("expected 200, got %d", data.Attempt.ResponseStatusCode)
	}
}

// ═══════════════════════════════════════════════════════════════════
// ERROR TYPES TESTS
// ═══════════════════════════════════════════════════════════════════

func TestAllErrorTypes(t *testing.T) {
	tests := []struct {
		code int
		want string
	}{
		{400, "*hooksniff.BadRequestError"},
		{401, "*hooksniff.UnauthorizedError"},
		{403, "*hooksniff.ForbiddenError"},
		{404, "*hooksniff.NotFoundError"},
		{409, "*hooksniff.ConflictError"},
		{422, "*hooksniff.UnprocessableEntityError"},
		{429, "*hooksniff.RateLimitError"},
		{500, "*hooksniff.InternalServerError"},
		{502, "*hooksniff.BadGatewayError"},
		{503, "*hooksniff.ServiceUnavailableError"},
		{504, "*hooksniff.GatewayTimeoutError"},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("status_%d", tt.code), func(t *testing.T) {
			err := CreateErrorFromStatus(tt.code, nil, nil)
			if err == nil {
				t.Fatal("expected error")
			}
			if got := fmt.Sprintf("%T", err); got != tt.want {
				t.Errorf("got %s, want %s", got, tt.want)
			}
		})
	}
}

func TestRateLimitErrorRetryAfter(t *testing.T) {
	headers := map[string][]string{"Retry-After": {"30"}}
	err := CreateErrorFromStatus(429, nil, headers)
	if err == nil {
		t.Fatal("expected error")
	}
	if rlErr, ok := err.(*RateLimitError); ok {
		if rlErr.RetryAfter != 30 {
			t.Errorf("expected retryAfter=30, got %d", rlErr.RetryAfter)
		}
	}
}

// ═══════════════════════════════════════════════════════════════════
// WEBHOOK EVENT BACKWARD COMPAT
// ═══════════════════════════════════════════════════════════════════

func TestWebhookEventGet(t *testing.T) {
	event := &WebhookEvent{
		Event:     "test",
		Data:      map[string]interface{}{"key": "value", "num": 42.0},
		Timestamp: "t",
	}
	if v := event.Get("key"); v != "value" {
		t.Errorf("expected value, got %v", v)
	}
	if v := event.Get("missing"); v != nil {
		t.Errorf("expected nil, got %v", v)
	}
}

func TestWebhookEventType(t *testing.T) {
	event := &WebhookEvent{Event: "endpoint.created"}
	if event.EventType() != "endpoint.created" {
		t.Errorf("expected endpoint.created, got %s", event.EventType())
	}
}

func TestWebhookEventFromJSON(t *testing.T) {
	jsonStr := `{"event":"endpoint.created","data":{"appId":"a1","endpointId":"e1"},"timestamp":"2026-05-19"}`
	var event WebhookEvent
	if err := json.Unmarshal([]byte(jsonStr), &event); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if event.Event != "endpoint.created" {
		t.Errorf("expected endpoint.created, got %s", event.Event)
	}
}

// ═══════════════════════════════════════════════════════════════════
// SDK VERSION
// ═══════════════════════════════════════════════════════════════════

func TestVersionNotEmpty(t *testing.T) {
	if Version == "" {
		t.Error("Version should not be empty")
	}
}

func TestVersionFormat(t *testing.T) {
	// Should be semver-like
	if len(Version) < 3 {
		t.Errorf("Version too short: %s", Version)
	}
}

// ═══════════════════════════════════════════════════════════════════
// EDGE CASES
// ═══════════════════════════════════════════════════════════════════

func TestParseEmptyPayload(t *testing.T) {
	wh := NewWebhook("whsec_dGVzdA==")
	payload := `{}`
	ts := nowUnix()
	sig := sign("whsec_dGVzdA==", "msg_1", ts, payload)

	event, err := wh.VerifyAndParse([]byte(payload), map[string]string{
		"webhook-id":        "msg_1",
		"webhook-timestamp": fmt.Sprintf("%d", ts),
		"webhook-signature": sig,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if event.Event != "" {
		t.Errorf("expected empty event, got %s", event.Event)
	}
}

func TestParseMinimalPayload(t *testing.T) {
	event := &WebhookEvent{
		Event:     "test",
		Data:      map[string]interface{}{},
		Timestamp: "",
	}
	if event.Event != "test" {
		t.Errorf("expected test, got %s", event.Event)
	}
}

func TestParseDataWithNestedJSON(t *testing.T) {
	jsonStr := `{"event":"endpoint.created","data":{"appId":"a1","endpointId":"e1","nested":{"key":"val"}},"timestamp":"t"}`
	var event WebhookEvent
	json.Unmarshal([]byte(jsonStr), &event)

	data, err := event.ParseEndpointCreatedData()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data.AppID != "a1" {
		t.Errorf("expected a1, got %s", data.AppID)
	}
}

// sign helper (duplicated from hooksniff_test.go for test independence)
func signForComprehensive(secret, msgID string, timestamp int64, payload string) string {
	decoded, _ := base64.StdEncoding.DecodeString(secret)
	toSign := fmt.Sprintf("%s.%d.%s", msgID, timestamp, payload)
	mac := hmac.New(sha256.New, decoded)
	mac.Write([]byte(toSign))
	sig := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	return "v1," + sig
}
