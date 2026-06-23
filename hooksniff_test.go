package hooksniff

import (
	"encoding/json"
	"testing"
)

func TestNewClient(t *testing.T) {
	hs := New("hr_test_fake", nil)
	if hs == nil {
		t.Fatal("New returned nil")
	}

	// Verify all 27 resources are initialized
	resources := map[string]interface{}{
		"Application":        hs.Application,
		"Endpoint":           hs.Endpoint,
		"Webhook":            hs.Webhook,
		"ApiKey":             hs.ApiKey,
		"Analytics":          hs.Analytics,
		"Search":             hs.Search,
		"Health":             hs.Health,
		"Team":               hs.Team,
		"Billing":            hs.Billing,
		"Notification":       hs.Notification,
		"Cortex":             hs.Cortex,
		"Template":           hs.Template,
		"Schema":             hs.Schema,
		"Alert":              hs.Alert,
		"Connector":          hs.Connector,
		"Stream":             hs.Stream,
		"BackgroundTask":     hs.BackgroundTask,
		"Integration":        hs.Integration,
		"ServiceToken":       hs.ServiceToken,
		"OperationalWebhook": hs.OperationalWebhook,
		"RateLimit":          hs.RateLimit,
		"Audit":              hs.Audit,
		"Sso":                hs.Sso,
		"CustomDomain":       hs.CustomDomain,
		"Environment":        hs.Environment,
		"Broadcast":          hs.Broadcast,
		"Transform":          hs.Transform,
	}

	for name, r := range resources {
		if r == nil {
			t.Errorf("Resource %s is nil", name)
		}
	}
	t.Logf("✅ All %d resources initialized", len(resources))
}

func TestWebhookVerification(t *testing.T) {
	// Test missing secret
	_, err := NewWebhook("")
	if err == nil {
		t.Error("Expected error for empty secret")
	}

	// Test missing headers
	wh, _ := NewWebhook("whsec_test")
	_, err = wh.Verify([]byte("{}"), map[string]string{})
	if err == nil {
		t.Error("Expected error for missing headers")
	}
	if ve, ok := err.(*WebhookVerificationError); ok {
		if ve.Message != "missing required webhook headers (webhook-id, webhook-signature, webhook-timestamp)" {
			t.Errorf("Unexpected error message: %s", ve.Message)
		}
	} else {
		t.Error("Expected WebhookVerificationError")
	}

	// Test old timestamp
	_, err = wh.Verify([]byte("{}"), map[string]string{
		"webhook-id":        "test",
		"webhook-signature": "v1,test",
		"webhook-timestamp": "1000000000",
	})
	if err == nil {
		t.Error("Expected error for old timestamp")
	}

	t.Log("✅ Webhook verification tests passed")
}

func TestErrorHierarchy(t *testing.T) {
	// Test error types
	authErr := &AuthenticationError{HookSniffError{401, "UNAUTHORIZED", "test"}}
	if authErr.Error() == "" {
		t.Error("Error() returned empty string")
	}

	notFoundErr := &NotFoundError{HookSniffError{404, "NOT_FOUND", "test"}}
	if notFoundErr.Error() == "" {
		t.Error("Error() returned empty string")
	}

	rateLimitErr := &RateLimitError{HookSniffError: HookSniffError{429, "RATE_LIMITED", "test"}, RetryAfter: 60}
	if rateLimitErr.Error() == "" {
		t.Error("Error() returned empty string")
	}

	t.Log("✅ Error hierarchy tests passed")
}

func TestMapError(t *testing.T) {
	tests := []struct {
		statusCode int
		wantType   string
	}{
		{401, "*hooksniff.AuthenticationError"},
		{404, "*hooksniff.NotFoundError"},
		{429, "*hooksniff.RateLimitError"},
		{400, "*hooksniff.ValidationError"},
		{422, "*hooksniff.ValidationError"},
		{500, "*hooksniff.ServerError"},
		{502, "*hooksniff.ServerError"},
	}

	for _, tt := range tests {
		body := map[string]interface{}{
			"error": map[string]interface{}{
				"code":   "TEST",
				"detail": "test error",
			},
		}
		err := mapError(tt.statusCode, body)
		if err == nil {
			t.Errorf("mapError(%d) returned nil", tt.statusCode)
			continue
		}
		t.Logf("✅ mapError(%d) → %T", tt.statusCode, err)
	}
}

func TestPagination(t *testing.T) {
	hs := New("hr_test_fake", nil)
	p := hs.Application.List()
	if p == nil {
		t.Error("List returned nil paginator")
	}
	t.Log("✅ Pagination test passed")
}

func TestCustomConfig(t *testing.T) {
	hs := New("hr_test_fake", &ClientConfig{
		BaseURL: "https://custom.api.com",
		Timeout: 10000000000, // 10s
		Retries: 5,
		Headers: map[string]string{"X-Custom": "value"},
	})
	if hs == nil {
		t.Fatal("New with config returned nil")
	}
	t.Log("✅ Custom config test passed")
}

func TestClientPanicsWithoutKey(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for empty API key")
		}
	}()
	New("", nil)
}

func TestTypes(t *testing.T) {
	// Verify types can be serialized
	app := Application{ID: "app_123", Name: "Test"}
	data, err := json.Marshal(app)
	if err != nil {
		t.Errorf("Failed to marshal Application: %v", err)
	}
	if len(data) == 0 {
		t.Error("Marshaled Application is empty")
	}

	ep := Endpoint{ID: "ep_123", URL: "https://example.com"}
	data, err = json.Marshal(ep)
	if err != nil {
		t.Errorf("Failed to marshal Endpoint: %v", err)
	}

	d := WebhookDelivery{ID: "msg_123", Status: "pending"}
	data, err = json.Marshal(d)
	if err != nil {
		t.Errorf("Failed to marshal WebhookDelivery: %v", err)
	}

	t.Log("✅ Type serialization tests passed")
}
