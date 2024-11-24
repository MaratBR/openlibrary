package server

import "net/http"

type httpError struct {
	StatusCode int
	Message    string
}

func httpErr(statusCode int, messages string) error {
	return &httpError{
		StatusCode: statusCode,
		Message:    messages,
	}
}

func (h *httpError) Error() string {
	if h.Message == "" {
		return "no message provided"
	}

	return h.Message
}

var (
	ErrHttpRequestBodyTooLarge = httpErr(http.StatusRequestEntityTooLarge, "request body too large")
)
