package retry

import (
	"errors"
	"fmt"
)

var (
	ErrCircuitBreaker  = errors.New("operation failed due to circuit breaker")
	ErrTemporary       = errors.New("temporary error")
	ErrNetwork         = errors.New("network failure")
	ErrTimeout         = errors.New("timeout occurred")
	ErrMaxRetries      = errors.New("operation failed after maximum retries")
	ErrContextCanceled = errors.New("operation canceled due to context timeout")
)

type HTTPError struct {
	StatusCode int
	Message    string
}

func (e HTTPError) Error() string {
	return fmt.Sprintf("HTTP %d: %s", e.StatusCode, e.Message)
}
