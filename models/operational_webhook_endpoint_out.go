package models

import "time"

type OperationalWebhookEndpointOut struct {
	Id        string    `json:"id"`
	Url       string    `json:"url"`
	Enabled   bool      `json:"enabled"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
