// Package hooksniff this file is @generated DO NOT EDIT
package models

type ListResponseEventTypeOut struct {
	Data         []EventTypeOut `json:"data"`
	Done         bool           `json:"done"`
	Iterator     *string        `json:"iterator,omitempty"`
	PrevIterator *string        `json:"prevIterator,omitempty"`
}

// Implement ListResponse interface for pagination
func (r *ListResponseEventTypeOut) GetData() []EventTypeOut { return r.Data }
func (r *ListResponseEventTypeOut) GetDone() bool { return r.Done }
func (r *ListResponseEventTypeOut) GetIterator() *string { return r.Iterator }
