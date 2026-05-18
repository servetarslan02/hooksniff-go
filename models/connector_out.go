package models

import "time"

type ConnectorOut struct {
	Id          string           `json:"id"`
	Name        string           `json:"name"`
	Description *string          `json:"description,omitempty"`
	ProductType *ConnectorProduct `json:"productType,omitempty"`
	Enabled     bool             `json:"enabled"`
	CreatedAt   time.Time        `json:"createdAt"`
	UpdatedAt   time.Time        `json:"updatedAt"`
}
