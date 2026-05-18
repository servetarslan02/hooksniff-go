package models

import "time"

type StreamSinkOut struct {
	Id        string    `json:"id"`
	StreamId  string    `json:"streamId"`
	Type      string    `json:"type"`
	Enabled   bool      `json:"enabled"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
