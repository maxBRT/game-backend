package client

import (
	"context"
	"errors"
	"net"
	"net/http"
	"syscall"
)

var (
	ErrInvalidTicketID    = errors.New("invalid or unknown ticket id")
	ErrPlayerFieldMissing = errors.New("player id , name or role field missing")
	ErrInvalidRole        = errors.New("role must be either 'survivor' or 'killer'")
)

type APIError struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}

func (e *APIError) Error() string {
	return e.Message
}

func IsRetryableError(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return true
	}

	if errors.Is(err, context.Canceled) {
		return false
	}

	var netErr net.Error
	if errors.As(err, &netErr) {
		if netErr.Timeout() {
			return true
		}

	}

	var opErr *net.OpError
	if errors.As(err, &opErr) {
		var syscallErr syscall.Errno
		if errors.As(opErr.Err, &syscallErr) {
			if syscallErr == syscall.ECONNREFUSED || syscallErr == syscall.ECONNRESET {
				return true
			}
		}
	}

	var APIErr *APIError

	if errors.As(err, &APIErr) {
		if APIErr.StatusCode >= 500 && APIErr.StatusCode <= 599 {
			if APIErr.StatusCode == http.StatusNotImplemented {
				return false
			}
		}

		return true
	}

	return false
}
