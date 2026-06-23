package hooksniff

import "encoding/json"

// HookSniff is the main SDK client.
type HookSniff struct {
	http *httpClient

	Application       *ApplicationResource
	Endpoint          *EndpointResource
	Webhook           *WebhookResource
	ApiKey            *ApiKeyResource
	Analytics         *AnalyticsResource
	Search            *SearchResource
	Health            *HealthResource
	Team              *TeamResource
	Billing           *BillingResource
	Notification      *NotificationResource
	Cortex            *CortexResource
	Template          *TemplateResource
	Schema            *SchemaResource
	Alert             *AlertResource
	Connector         *ConnectorResource
	Stream            *StreamResource
	BackgroundTask    *BackgroundTaskResource
	Integration       *IntegrationResource
	ServiceToken      *ServiceTokenResource
	OperationalWebhook *OperationalWebhookResource
	RateLimit         *RateLimitResource
	Audit             *AuditResource
	Sso               *SsoResource
	CustomDomain      *CustomDomainResource
	Environment       *EnvironmentResource
	Broadcast         *BroadcastResource
	Transform         *TransformResource
}

// New creates a new HookSniff client.
func New(apiKey string, cfg *ClientConfig) *HookSniff {
	h := newHTTPClient(apiKey, cfg)
	return &HookSniff{
		http:              h,
		Application:       &ApplicationResource{h},
		Endpoint:          &EndpointResource{h},
		Webhook:           &WebhookResource{h},
		ApiKey:            &ApiKeyResource{h},
		Analytics:         &AnalyticsResource{h},
		Search:            &SearchResource{h},
		Health:            &HealthResource{h},
		Team:              &TeamResource{h},
		Billing:           &BillingResource{h},
		Notification:      &NotificationResource{h},
		Cortex:            &CortexResource{h},
		Template:          &TemplateResource{h},
		Schema:            &SchemaResource{h},
		Alert:             &AlertResource{h},
		Connector:         &ConnectorResource{h},
		Stream:            &StreamResource{h},
		BackgroundTask:    &BackgroundTaskResource{h},
		Integration:       &IntegrationResource{h},
		ServiceToken:      &ServiceTokenResource{h},
		OperationalWebhook: &OperationalWebhookResource{h},
		RateLimit:         &RateLimitResource{h},
		Audit:             &AuditResource{h},
		Sso:               &SsoResource{h},
		CustomDomain:      &CustomDomainResource{h},
		Environment:       &EnvironmentResource{h},
		Broadcast:         &BroadcastResource{h},
		Transform:         &TransformResource{h},
	}
}

// Me returns the current user profile.
func (c *HookSniff) Me() (*User, error) {
	resp, err := c.http.request("GET", "/v1/auth/me", nil)
	if err != nil {
		return nil, err
	}
	var user User
	return &user, json.Unmarshal(resp, &user)
}
