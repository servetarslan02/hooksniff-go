package models

import "time"

type StreamEventTypeOut struct {
	Id          string    `json:"id"`
	Description *string   `json:"description,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
