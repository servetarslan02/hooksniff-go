package hooksniff

import (
	"net/http"
	"strconv"
)

// ResponseMetadata holds metadata from the last API response.
//
// Access via client.LastResponse after any API call.
type ResponseMetadata struct {
	// StatusCode is the HTTP status code
	StatusCode int
	// RequestID is the x-request-id header (for debugging with HookSniff support)
	RequestID string
	// RateLimitRemaining is the x-ratelimit-remaining header
	RateLimitRemaining *int
	// RateLimitReset is the x-ratelimit-reset header (Unix timestamp)
	RateLimitReset *int
	// Headers are all response headers
	Headers http.Header
}

// NewResponseMetadata creates ResponseMetadata from an http.Response.
func NewResponseMetadata(resp *http.Response) *ResponseMetadata {
	if resp == nil {
		return nil
	}

	rm := &ResponseMetadata{
		StatusCode: resp.StatusCode,
		RequestID:  resp.Header.Get("x-request-id"),
		Headers:    resp.Header,
	}

	if v := resp.Header.Get("x-ratelimit-remaining"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			rm.RateLimitRemaining = &n
		}
	}

	if v := resp.Header.Get("x-ratelimit-reset"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			rm.RateLimitReset = &n
		}
	}

	return rm
}
