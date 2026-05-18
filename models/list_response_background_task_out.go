// Package hooksniff this file is @generated DO NOT EDIT
package models

type ListResponseBackgroundTaskOut struct {
	Data         []BackgroundTaskOut `json:"data"`
	Done         bool                `json:"done"`
	Iterator     *string             `json:"iterator,omitempty"`
	PrevIterator *string             `json:"prevIterator,omitempty"`
}

// Implement ListResponse interface for pagination
func (r *ListResponseBackgroundTaskOut) GetData() []BackgroundTaskOut { return r.Data }
func (r *ListResponseBackgroundTaskOut) GetDone() bool { return r.Done }
func (r *ListResponseBackgroundTaskOut) GetIterator() *string { return r.Iterator }
