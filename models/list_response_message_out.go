// Package hooksniff this file is @generated DO NOT EDIT
package models

type ListResponseMessageOut struct {
	Data         []MessageOut `json:"data"`
	Done         bool         `json:"done"`
	Iterator     *string      `json:"iterator,omitempty"`
	PrevIterator *string      `json:"prevIterator,omitempty"`
}

// Implement ListResponse interface for pagination
func (r *ListResponseMessageOut) GetData() []MessageOut { return r.Data }
func (r *ListResponseMessageOut) GetDone() bool         { return r.Done }
func (r *ListResponseMessageOut) GetIterator() *string   { return r.Iterator }
