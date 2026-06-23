package hooksniff

import "fmt"

// HookSniffError is the base error type for all SDK errors.
type HookSniffError struct {
	StatusCode int    `json:"status_code"`
	Code       string `json:"code"`
	Detail     string `json:"detail"`
}

func (e *HookSniffError) Error() string {
	return fmt.Sprintf("HookSniff API error %d (%s): %s", e.StatusCode, e.Code, e.Detail)
}

// AuthenticationError indicates an invalid or missing API key.
type AuthenticationError struct{ HookSniffError }

// NotFoundError indicates the requested resource does not exist.
type NotFoundError struct{ HookSniffError }

// RateLimitError indicates the request was rate-limited.
type RateLimitError struct {
	HookSniffError
	RetryAfter int `json:"retry_after"`
}

// ValidationError indicates the request body failed validation.
type ValidationError struct{ HookSniffError }

// ServerError indicates an internal server error.
type ServerError struct{ HookSniffError }

func mapError(statusCode int, body map[string]interface{}) error {
	errObj, _ := body["error"].(map[string]interface{})
	code, _ := errObj["code"].(string)
	detail, _ := errObj["detail"].(string)
	if detail == "" {
		detail, _ = errObj["message"].(string)
	}
	if detail == "" {
		detail = "Unknown error"
	}
	if code == "" {
		code = "UNKNOWN"
	}

	base := HookSniffError{StatusCode: statusCode, Code: code, Detail: detail}

	switch statusCode {
	case 401:
		return &AuthenticationError{base}
	case 404:
		return &NotFoundError{base}
	case 429:
		return &RateLimitError{HookSniffError: base, RetryAfter: 60}
	case 400, 422:
		return &ValidationError{base}
	default:
		if statusCode >= 500 {
			return &ServerError{base}
		}
		return &base
	}
}
