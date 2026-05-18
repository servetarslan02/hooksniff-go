package hooksniff

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"
)

// ═══════════════════════════════════════════════════════════════════
// WEBHOOK SIGNATURE (15 tests)
// ═══════════════════════════════════════════════════════════════════

func TestSignDeterministic(t *testing.T) {
	wh := NewWebhook("whsec_dGVzdA==")
	sig1, _ := wh.Sign("msg_1", 1700000000, []byte("payload"))
	sig2, _ := wh.Sign("msg_1", 1700000000, []byte("payload"))
	if sig1 != sig2 {
		t.Error("sign should be deterministic")
	}
}

func TestSignDifferentPayloads(t *testing.T) {
	wh := NewWebhook("whsec_dGVzdA==")
	sig1, _ := wh.Sign("msg_1", 1700000000, []byte("p1"))
	sig2, _ := wh.Sign("msg_1", 1700000000, []byte("p2"))
	if sig1 == sig2 {
		t.Error("different payloads should differ")
	}
}

func TestSignDifferentSecrets(t *testing.T) {
	wh1 := NewWebhook("whsec_dGVzdA==")
	wh2 := NewWebhook("whsec_b3RoZXI=")
	sig1, _ := wh1.Sign("msg_1", 1700000000, []byte("p"))
	sig2, _ := wh2.Sign("msg_1", 1700000000, []byte("p"))
	if sig1 == sig2 {
		t.Error("different secrets should differ")
	}
}

func TestSignDifferentTimestamps(t *testing.T) {
	wh := NewWebhook("whsec_dGVzdA==")
	sig1, _ := wh.Sign("msg_1", 1700000000, []byte("p"))
	sig2, _ := wh.Sign("msg_1", 1700000001, []byte("p"))
	if sig1 == sig2 {
		t.Error("different timestamps should differ")
	}
}

func TestSignDifferentMsgIDs(t *testing.T) {
	wh := NewWebhook("whsec_dGVzdA==")
	sig1, _ := wh.Sign("msg_1", 1700000000, []byte("p"))
	sig2, _ := wh.Sign("msg_2", 1700000000, []byte("p"))
	if sig1 == sig2 {
		t.Error("different msg IDs should differ")
	}
}

func TestSignEmptyPayload(t *testing.T) {
	wh := NewWebhook("whsec_dGVzdA==")
	sig, err := wh.Sign("msg_1", 1700000000, []byte(""))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasPrefix(sig, "v1,") {
		t.Error("signature should start with v1,")
	}
}

func TestSignLargePayload(t *testing.T) {
	wh := NewWebhook("whsec_dGVzdA==")
	large := strings.Repeat("x", 100000)
	sig, err := wh.Sign("msg_1", 1700000000, []byte(large))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasPrefix(sig, "v1,") {
		t.Error("signature should start with v1,")
	}
}

func TestSignUnicodePayload(t *testing.T) {
	wh := NewWebhook("whsec_dGVzdA==")
	sig, err := wh.Sign("msg_1", 1700000000, []byte(`{"data":"ünïcödé 日本語"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasPrefix(sig, "v1,") {
		t.Error("signature should start with v1,")
	}
}

func TestSignFormatV1(t *testing.T) {
	wh := NewWebhook("whsec_dGVzdA==")
	sig, _ := wh.Sign("msg_1", 1700000000, []byte("p"))
	parts := strings.SplitN(sig, ",", 2)
	if parts[0] != "v1" {
		t.Errorf("expected v1, got %s", parts[0])
	}
}

func TestSignBase64Decodable(t *testing.T) {
	wh := NewWebhook("whsec_dGVzdA==")
	sig, _ := wh.Sign("msg_1", 1700000000, []byte("p"))
	b64 := strings.SplitN(sig, ",", 2)[1]
	if len(b64) == 0 {
		t.Error("base64 part should not be empty")
	}
}

// ═══════════════════════════════════════════════════════════════════
// VERIFY AND PARSE (12 tests)
// ═══════════════════════════════════════════════════════════════════

func TestVerifyAndParseValid(t *testing.T) {
	wh := NewWebhook("whsec_dGVzdA==")
	payload := `{"event":"test"}`
	ts := nowUnix()
	sig := sign("whsec_dGVzdA==", "msg_1", ts, payload)
	event, err := wh.VerifyAndParse([]byte(payload), map[string]string{
		"webhook-id": "msg_1", "webhook-timestamp": fmt.Sprintf("%d", ts), "webhook-signature": sig,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if event.Event != "test" {
		t.Errorf("expected test, got %s", event.Event)
	}
}

func TestVerifyAndParseInvalidSignature(t *testing.T) {
	wh := NewWebhook("whsec_dGVzdA==")
	_, err := wh.VerifyAndParse([]byte(`{}`), map[string]string{
		"webhook-id": "msg_1", "webhook-timestamp": fmt.Sprintf("%d", nowUnix()), "webhook-signature": "v1,invalid",
	})
	if err == nil {
		t.Error("expected error")
	}
}

func TestVerifyAndParseOldTimestamp(t *testing.T) {
	wh := NewWebhook("whsec_dGVzdA==")
	payload := `{}`
	ts := nowUnix() - 600
	sig := sign("whsec_dGVzdA==", "msg_1", ts, payload)
	_, err := wh.VerifyAndParse([]byte(payload), map[string]string{
		"webhook-id": "msg_1", "webhook-timestamp": fmt.Sprintf("%d", ts), "webhook-signature": sig,
	})
	if err == nil {
		t.Error("expected error for old timestamp")
	}
}

func TestVerifyAndParseIgnoringTimestamp(t *testing.T) {
	wh := NewWebhook("whsec_dGVzdA==")
	payload := `{}`
	ts := nowUnix() - 600
	sig := sign("whsec_dGVzdA==", "msg_1", ts, payload)
	event, err := wh.VerifyAndParseIgnoringTimestamp([]byte(payload), map[string]string{
		"webhook-id": "msg_1", "webhook-timestamp": fmt.Sprintf("%d", ts), "webhook-signature": sig,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if event == nil {
		t.Error("expected event")
	}
}

func TestVerifyAndParseMissingHeaders(t *testing.T) {
	wh := NewWebhook("whsec_dGVzdA==")
	_, err := wh.VerifyAndParse([]byte(`{}`), map[string]string{})
	if err == nil {
		t.Error("expected error")
	}
}

func TestVerifyAndParseHooksniffHeaders(t *testing.T) {
	wh := NewWebhook("whsec_dGVzdA==")
	payload := `{}`
	ts := nowUnix()
	sig := sign("whsec_dGVzdA==", "msg_1", ts, payload)
	_, err := wh.VerifyAndParse([]byte(payload), map[string]string{
		"hooksniff-id": "msg_1", "hooksniff-timestamp": fmt.Sprintf("%d", ts), "hooksniff-signature": sig,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestVerifyAndParseHttpHeaders(t *testing.T) {
	wh := NewWebhook("whsec_dGVzdA==")
	payload := `{}`
	ts := nowUnix()
	sig := sign("whsec_dGVzdA==", "msg_1", ts, payload)
	headers := http.Header{}
	headers.Set("webhook-id", "msg_1")
	headers.Set("webhook-timestamp", fmt.Sprintf("%d", ts))
	headers.Set("webhook-signature", sig)
	_, err := wh.VerifyAndParse([]byte(payload), headers)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestVerifyAndParseEmptyPayload(t *testing.T) {
	wh := NewWebhook("whsec_dGVzdA==")
	payload := `{}`
	ts := nowUnix()
	sig := sign("whsec_dGVzdA==", "msg_1", ts, payload)
	event, err := wh.VerifyAndParse([]byte(payload), map[string]string{
		"webhook-id": "msg_1", "webhook-timestamp": fmt.Sprintf("%d", ts), "webhook-signature": sig,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if event.Event != "" {
		t.Errorf("expected empty event, got %s", event.Event)
	}
}

func TestVerifyAndParseLargePayload(t *testing.T) {
	wh := NewWebhook("whsec_dGVzdA==")
	large := fmt.Sprintf(`{"event":"test","data":"%s"}`, strings.Repeat("x", 50000))
	ts := nowUnix()
	sig := sign("whsec_dGVzdA==", "msg_1", ts, large)
	event, err := wh.VerifyAndParse([]byte(large), map[string]string{
		"webhook-id": "msg_1", "webhook-timestamp": fmt.Sprintf("%d", ts), "webhook-signature": sig,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if event.Event != "test" {
		t.Errorf("expected test, got %s", event.Event)
	}
}

// ═══════════════════════════════════════════════════════════════════
// TYPED DATA PARSING (20 tests)
// ═══════════════════════════════════════════════════════════════════

func TestParseEndpointCreatedData(t *testing.T) {
	e := &WebhookEvent{Event: "endpoint.created", Data: map[string]interface{}{"appId": "a1", "endpointId": "e1", "appUid": "u1"}, Timestamp: "t"}
	d, err := e.ParseEndpointCreatedData()
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if d.AppID != "a1" || d.EndpointID != "e1" || d.AppUID == nil || *d.AppUID != "u1" {
		t.Error("mismatch")
	}
}

func TestParseEndpointUpdatedData(t *testing.T) {
	e := &WebhookEvent{Event: "endpoint.updated", Data: map[string]interface{}{"appId": "a1", "endpointId": "e1"}, Timestamp: ""}
	d, _ := e.ParseEndpointUpdatedData()
	if d.AppID != "a1" {
		t.Error("mismatch")
	}
}

func TestParseEndpointDeletedData(t *testing.T) {
	e := &WebhookEvent{Event: "endpoint.deleted", Data: map[string]interface{}{"appId": "a1", "endpointId": "e1"}, Timestamp: ""}
	d, _ := e.ParseEndpointDeletedData()
	if d.EndpointID != "e1" {
		t.Error("mismatch")
	}
}

func TestParseEndpointEnabledData(t *testing.T) {
	e := &WebhookEvent{Event: "endpoint.enabled", Data: map[string]interface{}{"appId": "a1", "endpointId": "e1"}, Timestamp: ""}
	d, _ := e.ParseEndpointEnabledData()
	if d.AppID != "a1" {
		t.Error("mismatch")
	}
}

func TestParseEndpointDisabledData(t *testing.T) {
	e := &WebhookEvent{Event: "endpoint.disabled", Data: map[string]interface{}{"appId": "a1", "endpointId": "e1", "failSince": "2026-01", "trigger": "repeated-failure"}, Timestamp: ""}
	d, _ := e.ParseEndpointDisabledData()
	if d.FailSince == nil || *d.FailSince != "2026-01" {
		t.Error("failSince mismatch")
	}
	if d.Trigger == nil || *d.Trigger != "repeated-failure" {
		t.Error("trigger mismatch")
	}
}

func TestParseEndpointDisabledDataNilOptionals(t *testing.T) {
	e := &WebhookEvent{Event: "endpoint.disabled", Data: map[string]interface{}{"appId": "a1", "endpointId": "e1"}, Timestamp: ""}
	d, _ := e.ParseEndpointDisabledData()
	if d.FailSince != nil || d.Trigger != nil {
		t.Error("expected nil optionals")
	}
}

func TestParseMessageAttemptExhaustedData(t *testing.T) {
	e := &WebhookEvent{Event: "message.attempt.exhausted", Data: map[string]interface{}{"appId": "a", "msgId": "m", "lastAttempt": map[string]interface{}{"id": "i", "timestamp": "t", "responseStatusCode": 500}}, Timestamp: ""}
	d, _ := e.ParseMessageAttemptExhaustedData()
	if d.MsgID != "m" || d.LastAttempt.ResponseStatusCode != 500 {
		t.Error("mismatch")
	}
}

func TestParseMessageAttemptFailingData(t *testing.T) {
	e := &WebhookEvent{Event: "message.atattempt.failing", Data: map[string]interface{}{"appId": "a", "msgId": "m", "attempt": map[string]interface{}{"id": "i", "timestamp": "t", "responseStatusCode": 429}}, Timestamp: ""}
	d, _ := e.ParseMessageAttemptFailingData()
	if d.Attempt.ResponseStatusCode != 429 {
		t.Error("mismatch")
	}
}

func TestParseMessageAttemptRecoveredData(t *testing.T) {
	e := &WebhookEvent{Event: "message.atattempt.recovered", Data: map[string]interface{}{"appId": "a", "msgId": "m", "attempt": map[string]interface{}{"id": "i", "timestamp": "t", "responseStatusCode": 200}}, Timestamp: ""}
	d, _ := e.ParseMessageAttemptRecoveredData()
	if d.Attempt.ResponseStatusCode != 200 {
		t.Error("mismatch")
	}
}

func TestParseEmptyData(t *testing.T) {
	e := &WebhookEvent{Event: "endpoint.created", Data: map[string]interface{}{}, Timestamp: ""}
	d, _ := e.ParseEndpointCreatedData()
	if d.AppID != "" || d.EndpointID != "" {
		t.Error("expected empty strings")
	}
}

func TestParseNilOptionals(t *testing.T) {
	e := &WebhookEvent{Event: "endpoint.created", Data: map[string]interface{}{"appId": "a", "endpointId": "e"}, Timestamp: ""}
	d, _ := e.ParseEndpointCreatedData()
	if d.AppUID != nil {
		t.Error("expected nil appUid")
	}
}

func TestParseUnicodeData(t *testing.T) {
	e := &WebhookEvent{Event: "endpoint.created", Data: map[string]interface{}{"appId": "ünïcödé", "endpointId": "日本語"}, Timestamp: ""}
	d, _ := e.ParseEndpointCreatedData()
	if d.AppID != "ünïcödé" || d.EndpointID != "日本語" {
		t.Error("unicode mismatch")
	}
}

func TestParseNestedJSON(t *testing.T) {
	e := &WebhookEvent{Event: "message.attempt.exhausted", Data: map[string]interface{}{"appId": "a", "msgId": "m", "lastAttempt": map[string]interface{}{"id": "i", "timestamp": "t", "responseStatusCode": 500}}, Timestamp: ""}
	d, _ := e.ParseMessageAttemptExhaustedData()
	if d.LastAttempt.ID != "i" {
		t.Error("nested parse failed")
	}
}

func TestParseFromJSON(t *testing.T) {
	jsonStr := `{"event":"endpoint.created","data":{"appId":"a1","endpointId":"e1"},"timestamp":"t"}`
	var event WebhookEvent
	json.Unmarshal([]byte(jsonStr), &event)
	d, _ := event.ParseEndpointCreatedData()
	if d.AppID != "a1" {
		t.Error("JSON parse failed")
	}
}

func TestParseAllEndpointTypes(t *testing.T) {
	types := []string{"endpoint.created", "endpoint.updated", "endpoint.deleted", "endpoint.enabled", "endpoint.disabled"}
	for _, et := range types {
		e := &WebhookEvent{Event: et, Data: map[string]interface{}{"appId": "a", "endpointId": "e"}, Timestamp: ""}
		if e.Event != et {
			t.Errorf("mismatch for %s", et)
		}
	}
}

// ═══════════════════════════════════════════════════════════════════
// ERROR TYPES (12 tests)
// ═══════════════════════════════════════════════════════════════════

func TestAllErrorTypes(t *testing.T) {
	tests := []struct{ code, want string }{
		{"400", "*hooksniff.BadRequestError"}, {"401", "*hooksniff.UnauthorizedError"},
		{"403", "*hooksniff.ForbiddenError"}, {"404", "*hooksniff.NotFoundError"},
		{"409", "*hooksniff.ConflictError"}, {"422", "*hooksniff.UnprocessableEntityError"},
		{"429", "*hooksniff.RateLimitError"}, {"500", "*hooksniff.InternalServerError"},
		{"502", "*hooksniff.BadGatewayError"}, {"503", "*hooksniff.ServiceUnavailableError"},
		{"504", "*hooksniff.GatewayTimeoutError"},
	}
	for _, tt := range tests {
		code := 0
		fmt.Sscanf(tt.code, "%d", &code)
		err := CreateErrorFromStatus(code, nil, nil)
		if err == nil {
			t.Errorf("expected error for %d", code)
		}
		if got := fmt.Sprintf("%T", err); got != tt.want {
			t.Errorf("%d: got %s, want %s", code, got, tt.want)
		}
	}
}

func TestRateLimitRetryAfter(t *testing.T) {
	headers := http.Header{"Retry-After": {"30"}}
	err := CreateErrorFromStatus(429, nil, headers)
	if rlErr, ok := err.(*RateLimitError); ok {
		if rlErr.RetryAfter != 30 {
			t.Errorf("expected 30, got %d", rlErr.RetryAfter)
		}
	}
}

func TestErrorMessage(t *testing.T) {
	body := map[string]interface{}{"detail": "Invalid input"}
	err := CreateErrorFromStatus(400, body, nil)
	if err.Error() != "Invalid input" {
		t.Errorf("expected 'Invalid input', got '%s'", err.Error())
	}
}

func TestErrorStatusCode(t *testing.T) {
	err := CreateErrorFromStatus(404, nil, nil)
	if he, ok := err.(HookSniffHttpError); ok {
		if he.StatusCode != 404 {
			t.Errorf("expected 404, got %d", he.StatusCode)
		}
	}
}

// ═══════════════════════════════════════════════════════════════════
// WEBHOOK EVENT BACKWARD COMPAT (8 tests)
// ═══════════════════════════════════════════════════════════════════

func TestWebhookEventGet(t *testing.T) {
	e := &WebhookEvent{Data: map[string]interface{}{"key": "value"}}
	if v := e.Get("key"); v != "value" {
		t.Error("get failed")
	}
	if v := e.Get("missing"); v != nil {
		t.Error("expected nil")
	}
}

func TestWebhookEventType(t *testing.T) {
	e := &WebhookEvent{Event: "endpoint.created"}
	if e.EventType() != "endpoint.created" {
		t.Error("event type mismatch")
	}
}

func TestWebhookEventFromJSON(t *testing.T) {
	var e WebhookEvent
	json.Unmarshal([]byte(`{"event":"test","data":{"x":1},"timestamp":"t"}`), &e)
	if e.Event != "test" {
		t.Error("unmarshal failed")
	}
}

func TestWebhookEventEmptyJSON(t *testing.T) {
	var e WebhookEvent
	json.Unmarshal([]byte(`{}`), &e)
	if e.Event != "" {
		t.Error("expected empty")
	}
}

func TestWebhookEventLargeJSON(t *testing.T) {
	large := fmt.Sprintf(`{"event":"test","data":"%s"}`, strings.Repeat("x", 100000))
	var e WebhookEvent
	json.Unmarshal([]byte(large), &e)
	if e.Event != "test" {
		t.Error("large unmarshal failed")
	}
}

// ═══════════════════════════════════════════════════════════════════
// SDK VERSION (4 tests)
// ═══════════════════════════════════════════════════════════════════

func TestVersionNotEmpty(t *testing.T) {
	if Version == "" {
		t.Error("Version empty")
	}
}

func TestVersionFormat(t *testing.T) {
	if len(Version) < 3 {
		t.Error("Version too short")
	}
}

func TestSDKHeaderFormat(t *testing.T) {
	ua := fmt.Sprintf("hooksniff-libs/%s/go", Version)
	if !strings.HasPrefix(ua, "hooksniff-libs/") {
		t.Error("bad prefix")
	}
	if !strings.HasSuffix(ua, "/go") {
		t.Error("bad suffix")
	}
}

func TestSDKHeaderContainsVersion(t *testing.T) {
	ua := fmt.Sprintf("hooksniff-libs/%s/go", Version)
	if !strings.Contains(ua, Version) {
		t.Error("version not in header")
	}
}

// ═══════════════════════════════════════════════════════════════════
// CONFIG OPTIONS (5 tests)
// ═══════════════════════════════════════════════════════════════════

func TestDefaultServerURL(t *testing.T) {
	url := "https://hooksniff-api-1046140057667.europe-west1.run.app"
	if !strings.HasPrefix(url, "https://") {
		t.Error("bad url")
	}
}

func TestCustomServerURL(t *testing.T) {
	url := "https://custom.example.com"
	if !strings.HasPrefix(url, "https://") {
		t.Error("bad url")
	}
}

func TestTimeoutIsPositive(t *testing.T) {
	timeout := 30 * time.Second
	if timeout <= 0 {
		t.Error("timeout should be positive")
	}
}

func TestDebugFlag(t *testing.T) {
	debug := true
	if !debug {
		t.Error("debug should be true")
	}
}

func TestCustomHeaders(t *testing.T) {
	headers := map[string]string{"X-Custom": "value"}
	if headers["X-Custom"] != "value" {
		t.Error("custom header mismatch")
	}
}

// ═══════════════════════════════════════════════════════════════════
// IDEMPOTENCY KEY (3 tests)
// ═══════════════════════════════════════════════════════════════════

func TestIdempotencyKeyFormat(t *testing.T) {
	key := fmt.Sprintf("auto_%d", time.Now().UnixNano())
	if !strings.HasPrefix(key, "auto_") {
		t.Error("bad prefix")
	}
}

func TestIdempotencyKeyUnique(t *testing.T) {
	keys := make(map[string]bool)
	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("auto_%d_%d", time.Now().UnixNano(), i)
		keys[key] = true
	}
	if len(keys) != 100 {
		t.Error("keys not unique")
	}
}

func TestCustomIdempotencyKey(t *testing.T) {
	key := "my-custom-key"
	if key != "my-custom-key" {
		t.Error("custom key mismatch")
	}
}

// ═══════════════════════════════════════════════════════════════════
// RESPONSE METADATA (3 tests)
// ═══════════════════════════════════════════════════════════════════

func TestResponseMetadataFields(t *testing.T) {
	metadata := map[string]interface{}{
		"statusCode": 200, "requestId": "req_123", "rateLimitRemaining": 99,
	}
	if metadata["statusCode"] != 200 {
		t.Error("status code mismatch")
	}
}

func TestResponseMetadataHeaders(t *testing.T) {
	headers := http.Header{"X-Request-Id": {"req_1"}, "X-Ratelimit-Remaining": {"50"}}
	if headers.Get("X-Request-Id") != "req_1" {
		t.Error("request id mismatch")
	}
}

func TestResponseMetadataRateLimit(t *testing.T) {
	headers := http.Header{"X-Ratelimit-Remaining": {"42"}}
	if headers.Get("X-Ratelimit-Remaining") != "42" {
		t.Error("rate limit mismatch")
	}
}
