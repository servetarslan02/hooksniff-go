// Package hooksniff this file is @generated DO NOT EDIT
package models

type ListResponseOperationalWebhookEndpointOut struct {
	Data         []OperationalWebhookEndpointOut `json:"data"`
	Done         bool                            `json:"done"`
	Iterator     *string                         `json:"iterator,omitempty"`
	PrevIterator *string                         `json:"prevIterator,omitempty"`
}

// Implement ListResponse interface for pagination
func (r *ListResponseOperationalWebhookEndpointOut) GetData() []OperationalWebhookEndpointOut { return r.Data }
func (r *ListResponseOperationalWebhookEndpointOut) GetDone() bool { return r.Done }
func (r *ListResponseOperationalWebhookEndpointOut) GetIterator() *string { return r.Iterator }
