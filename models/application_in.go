package models

type ApplicationIn struct {
	Name        string            `json:"name"`
	Metadata    map[string]string `json:"metadata,omitempty"`
	RateLimit   *uint16           `json:"rateLimit,omitempty"`
	ThrottleRate *uint16          `json:"throttleRate,omitempty"`
	Uid         *string           `json:"uid,omitempty"`
}
