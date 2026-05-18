package models

type PollingEndpointOut struct {
	Data    []PollingEndpointMessageOut `json:"data"`
	Done    bool                        `json:"done"`
	Iterator *string                    `json:"iterator,omitempty"`
}
