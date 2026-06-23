package hooksniff

import "time"

// Application represents a HookSniff application.
type Application struct {
	ID            string    `json:"id"`
	CustomerID    string    `json:"customer_id"`
	Name          string    `json:"name"`
	Description   *string   `json:"description"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	EndpointCount int       `json:"endpoint_count"`
}

// ApplicationCreate is the request body for creating an application.
type ApplicationCreate struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
}

// ApplicationUpdate is the request body for updating an application.
type ApplicationUpdate struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

// Endpoint represents a webhook endpoint.
type Endpoint struct {
	ID                  string                 `json:"id"`
	URL                 string                 `json:"url"`
	Description         *string                `json:"description"`
	IsActive            bool                   `json:"is_active"`
	RetryPolicy         map[string]interface{} `json:"retry_policy"`
	CreatedAt           time.Time              `json:"created_at"`
	AllowedIPs          []string               `json:"allowed_ips"`
	EventFilter         []string               `json:"event_filter"`
	CustomHeaders       map[string]string      `json:"custom_headers"`
	RoutingStrategy     string                 `json:"routing_strategy"`
	FallbackURL         *string                `json:"fallback_url"`
	AvgResponseMs       float64                `json:"avg_response_ms"`
	FailureStreak       int                    `json:"failure_streak"`
	Format              string                 `json:"format"`
	ApplicationID       string                 `json:"application_id"`
}

// EndpointCreate is the request body for creating an endpoint.
type EndpointCreate struct {
	URL             string            `json:"url"`
	ApplicationID   string            `json:"application_id"`
	Description     *string           `json:"description,omitempty"`
	AllowedIPs      []string          `json:"allowed_ips,omitempty"`
	EventFilter     []string          `json:"event_filter,omitempty"`
	CustomHeaders   map[string]string `json:"custom_headers,omitempty"`
	RoutingStrategy *string           `json:"routing_strategy,omitempty"`
	FallbackURL     *string           `json:"fallback_url,omitempty"`
}

// EndpointUpdate is the request body for updating an endpoint.
type EndpointUpdate struct {
	URL             *string           `json:"url,omitempty"`
	Description     *string           `json:"description,omitempty"`
	IsActive        *bool             `json:"is_active,omitempty"`
	AllowedIPs      []string          `json:"allowed_ips,omitempty"`
	EventFilter     []string          `json:"event_filter,omitempty"`
	CustomHeaders   map[string]string `json:"custom_headers,omitempty"`
	RoutingStrategy *string           `json:"routing_strategy,omitempty"`
	FallbackURL     *string           `json:"fallback_url,omitempty"`
}

// WebhookDelivery represents a webhook delivery record.
type WebhookDelivery struct {
	ID           string    `json:"id"`
	EndpointID   string    `json:"endpoint_id"`
	Event        string    `json:"event"`
	Status       string    `json:"status"`
	AttemptCount int       `json:"attempt_count"`
	ResponseStatus *int    `json:"response_status"`
	ReplayCount  int       `json:"replay_count"`
	CreatedAt    time.Time `json:"created_at"`
	IsTest       bool      `json:"is_test"`
}

// WebhookSend is the request body for sending a webhook.
type WebhookSend struct {
	EndpointID string                 `json:"endpoint_id"`
	Event      string                 `json:"event"`
	Data       map[string]interface{} `json:"data"`
	IsTest     *bool                  `json:"is_test,omitempty"`
}

// ApiKey represents an API key.
type ApiKey struct {
	ID           string     `json:"id"`
	Name         string     `json:"name"`
	APIKeyPrefix string     `json:"api_key_prefix"`
	CreatedAt    time.Time  `json:"created_at"`
	LastUsedAt   *time.Time `json:"last_used_at"`
	IsActive     bool       `json:"is_active"`
}

// ApiKeyCreated is returned when creating a new API key.
type ApiKeyCreated struct {
	ID      string `json:"id"`
	Key     string `json:"key"`
	Prefix  string `json:"prefix"`
	Message string `json:"message"`
}

// User represents the current user profile.
type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	Name         string    `json:"name"`
	Plan         string    `json:"plan"`
	WebhookLimit int       `json:"webhook_limit"`
	WebhookCount int       `json:"webhook_count"`
	IsAdmin      bool      `json:"is_admin"`
	CreatedAt    time.Time `json:"created_at"`
}

// SearchResult represents search results.
type SearchResult struct {
	Deliveries []WebhookDelivery `json:"deliveries"`
	Total      int               `json:"total"`
	Page       int               `json:"page"`
	PerPage    int               `json:"per_page"`
	Query      string            `json:"query"`
}

// HealthResponse represents the health check response.
type HealthResponse struct {
	Status   string                       `json:"status"`
	API      *HealthComponent             `json:"api"`
	Database *HealthComponent             `json:"database"`
	Redis    *HealthComponent             `json:"redis"`
	Queue    *HealthQueue                 `json:"queue"`
	Checks   map[string]*HealthComponent  `json:"checks"`
}

// HealthComponent represents a single health check component.
type HealthComponent struct {
	Status    string `json:"status"`
	LatencyMs *int64 `json:"latency_ms"`
	Error     string `json:"error,omitempty"`
}

// HealthQueue represents queue health metrics.
type HealthQueue struct {
	Pending    int `json:"pending"`
	Processing int `json:"processing"`
	Failed     int `json:"failed"`
}

// OutboundIPs represents the outbound IP addresses.
type OutboundIPs struct {
	IPs       []string `json:"ips"`
	UpdatedAt string   `json:"updated_at"`
}

// SecretRotateResponse is returned when rotating an endpoint secret.
type SecretRotateResponse struct {
	ID                 string `json:"id"`
	Message            string `json:"message"`
	OldSecretValidUntil string `json:"old_secret_valid_until"`
	SigningSecret       string `json:"signing_secret"`
}

// Subscription represents a billing subscription.
type Subscription struct {
	Plan           string  `json:"plan"`
	Status         string  `json:"status"`
	PaymentProvider *string `json:"payment_provider"`
	WebhookLimit   int     `json:"webhook_limit"`
	EndpointLimit  *int    `json:"endpoint_limit"`
	RetentionDays  *int    `json:"retention_days"`
}


