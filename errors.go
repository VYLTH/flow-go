package flow

import "fmt"

// FlowError is the base error type for all API errors.
type FlowError struct {
	Message string `json:"message"`
	Code    string `json:"code"`
	Status  int    `json:"status"`
}

func (e *FlowError) Error() string {
	return fmt.Sprintf("flow: %s (code=%s, status=%d)", e.Message, e.Code, e.Status)
}

// AuthenticationError indicates an invalid or missing API key.
type AuthenticationError struct{ FlowError }

// NotFoundError indicates the requested resource was not found.
type NotFoundError struct{ FlowError }

// ValidationError indicates invalid request parameters.
type ValidationError struct {
	FlowError
	Errors []map[string]string `json:"errors"`
}

// RateLimitError indicates too many requests.
type RateLimitError struct {
	FlowError
	RetryAfter float64 `json:"retry_after"`
}

// ServerError indicates a server-side error.
type ServerError struct{ FlowError }

// InvalidSignatureError indicates webhook signature verification failed.
type InvalidSignatureError struct {
	Message string
}

func (e *InvalidSignatureError) Error() string {
	if e.Message == "" {
		return "flow: invalid webhook signature"
	}
	return "flow: " + e.Message
}
