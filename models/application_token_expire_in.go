package models

type ApplicationTokenExpireIn struct {
	Expiry *int64 `json:"expiry,omitempty"`
}
