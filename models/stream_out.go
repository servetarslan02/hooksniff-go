package models

import "time"

type StreamOut struct {
	Id        string            `json:"id"`
	Name      *string           `json:"name,omitempty"`
	Metadata  map[string]string `json:"metadata"`
	Uid       *string           `json:"uid,omitempty"`
	CreatedAt time.Time         `json:"createdAt"`
	UpdatedAt time.Time         `json:"updatedAt"`
}
