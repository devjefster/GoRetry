package retry

import "time"

// / BackoffStrategy defines the type of backoff to apply.
type BackoffStrategy int

const (
	FixedBackoff BackoffStrategy = iota
	LinearBackoff
	ExponentialBackoff
	CustomBackoff
)

// Config struct now includes BackoffStrategy and CustomBackoffFunc
type Config struct {
	MaxRetryTimeout         time.Duration
	MaxRetries              int
	InitialBackoff          time.Duration
	MaxBackoff              time.Duration
	BackoffFactor           float64
	JitterFactor            float64
	RetryableErrors         []error
	BackoffStrategy         BackoffStrategy
	RetryableStatusCodes    []int
	CircuitBreakerThreshold int
	CustomBackoffFunc       func(attempt int) time.Duration // User-defined backoff function
	SuccessReset            int
}

// DefaultConfig provides default retry settings
var DefaultConfig = Config{
	MaxRetries:     5,
	InitialBackoff: 100 * time.Millisecond,
	MaxBackoff:     5 * time.Second,
	BackoffFactor:  2.0,
}

type CircuitBreaker struct {
	failures     int
	successes    int
	Threshold    int
	SuccessReset int // Number of consecutive successes before reset
}

func (cb *CircuitBreaker) ShouldRetry() bool {
	return cb.failures < cb.Threshold
}

func (cb *CircuitBreaker) RecordFailure() {
	cb.failures++
	cb.successes = 0
}

func (cb *CircuitBreaker) RecordSuccess() {
	cb.successes++
	if cb.successes >= cb.SuccessReset {
		cb.Reset()
	}
}

func (cb *CircuitBreaker) Reset() {
	cb.failures = 0
	cb.successes = 0
}
