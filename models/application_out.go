package models

import "time"

type ApplicationOut struct {
	CreatedAt   time.Time         `json:"createdAt"`
	Id          string            `json:"id"`
	Metadata    map[string]string `json:"metadata"`
	Name        string            `json:"name"`
	RateLimit   *uint16           `json:"rateLimit,omitempty"`
	ThrottleRate *uint16          `json:"throttleRate,omitempty"`
	Uid         *string           `json:"uid,omitempty"`
	UpdatedAt   time.Time         `json:"updatedAt"`
}
