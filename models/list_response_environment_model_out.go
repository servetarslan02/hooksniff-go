// Package hooksniff this file is @generated DO NOT EDIT
package models

type ListResponseEnvironmentModelOut struct {
	Data         []EnvironmentModelOut `json:"data"`
	Done         bool                  `json:"done"`
	Iterator     *string               `json:"iterator,omitempty"`
	PrevIterator *string               `json:"prevIterator,omitempty"`
}

// Implement ListResponse interface for pagination
func (r *ListResponseEnvironmentModelOut) GetData() []EnvironmentModelOut { return r.Data }
func (r *ListResponseEnvironmentModelOut) GetDone() bool { return r.Done }
func (r *ListResponseEnvironmentModelOut) GetIterator() *string { return r.Iterator }
