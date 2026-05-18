package models

type OperationalWebhookEndpointUpdate struct {
	Url     string            `json:"url"`
	Enabled *bool             `json:"enabled,omitempty"`
	FilterTypes []string      `json:"filterTypes,omitempty"`
	Description *string        `json:"description,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}
