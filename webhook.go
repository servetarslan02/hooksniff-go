package hooksniff

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// WebhookVerificationError indicates webhook signature verification failed.
type WebhookVerificationError struct {
	Message string
}

func (e *WebhookVerificationError) Error() string {
	return fmt.Sprintf("webhook verification failed: %s", e.Message)
}

// Webhook verifies incoming webhook payloads using Standard Webhooks HMAC-SHA256.
type Webhook struct {
	secret string
}

// NewWebhook creates a new Webhook verifier with the given secret (whsec_...).
func NewWebhook(secret string) (*Webhook, error) {
	if secret == "" {
		return nil, fmt.Errorf("webhook secret is required")
	}
	return &Webhook{secret: secret}, nil
}

// Verify verifies the webhook payload and returns the parsed event.
// Headers must contain webhook-id, webhook-signature, and webhook-timestamp.
func (w *Webhook) Verify(payload []byte, headers map[string]string) (json.RawMessage, error) {
	// Normalize headers to lowercase
	normalized := make(map[string]string, len(headers))
	for k, v := range headers {
		normalized[strings.ToLower(k)] = v
	}

	msgID := normalized["webhook-id"]
	msgSignature := normalized["webhook-signature"]
	msgTimestamp := normalized["webhook-timestamp"]

	if msgID == "" || msgSignature == "" || msgTimestamp == "" {
		return nil, &WebhookVerificationError{
			Message: "missing required webhook headers (webhook-id, webhook-signature, webhook-timestamp)",
		}
	}

	// Validate timestamp (reject if older than 5 minutes)
	timestamp, err := strconv.ParseInt(msgTimestamp, 10, 64)
	if err != nil {
		return nil, &WebhookVerificationError{Message: "invalid webhook timestamp"}
	}

	now := time.Now().Unix()
	if abs(now-timestamp) > 300 {
		return nil, &WebhookVerificationError{Message: "webhook timestamp is too old"}
	}

	// Compute expected signature
	toSign := fmt.Sprintf("%s.%s.%s", msgID, msgTimestamp, string(payload))
	expectedSignature := w.sign(toSign)

	// Check signatures (space-separated, format: "v1,<base64>")
	signatures := strings.Split(msgSignature, " ")
	valid := false
	for _, sig := range signatures {
		parts := strings.SplitN(sig, ",", 2)
		if len(parts) != 2 {
			continue
		}
		if parts[0] != "v1" {
			continue
		}
		if hmac.Equal([]byte(parts[1]), []byte(expectedSignature)) {
			valid = true
			break
		}
	}

	if !valid {
		return nil, &WebhookVerificationError{Message: "invalid webhook signature"}
	}

	return json.RawMessage(payload), nil
}

func (w *Webhook) sign(content string) string {
	secretStr := strings.TrimPrefix(w.secret, "whsec_")
	secretBytes, err := base64.StdEncoding.DecodeString(secretStr)
	if err != nil {
		secretBytes = []byte(secretStr)
	}
	mac := hmac.New(sha256.New, secretBytes)
	mac.Write([]byte(content))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}
