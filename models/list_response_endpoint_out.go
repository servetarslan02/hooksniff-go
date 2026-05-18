// Package hooksniff this file is @generated DO NOT EDIT
package models

type ListResponseEndpointOut struct {
	Data         []EndpointOut `json:"data"`
	Done         bool          `json:"done"`
	Iterator     *string       `json:"iterator,omitempty"`
	PrevIterator *string       `json:"prevIterator,omitempty"`
}

// Implement ListResponse interface for pagination
func (r *ListResponseEndpointOut) GetData() []EndpointOut { return r.Data }
func (r *ListResponseEndpointOut) GetDone() bool { return r.Done }
func (r *ListResponseEndpointOut) GetIterator() *string { return r.Iterator }
