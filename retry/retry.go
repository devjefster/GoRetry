package retry

import (
	"context"
	"errors"
	"net/http"
	"time"
)

import "log"

func Retry(ctx context.Context, fn func() error, cfg Config) error {
	if cfg.MaxRetries < 0 {
		cfg.MaxRetries = DefaultConfig.MaxRetries
	}
	if cfg.MaxRetryTimeout < 0 {
		cfg.MaxRetryTimeout = DefaultConfig.MaxRetryTimeout
	}

	cb := CircuitBreaker{Threshold: cfg.CircuitBreakerThreshold}
	var err error
	log.Printf("Starting retry process: MaxRetries=%d, BackoffStrategy=%v", cfg.MaxRetries, cfg.BackoffStrategy)

	for i := 0; i < cfg.MaxRetries; i++ {

		if !cb.ShouldRetry() {
			log.Println("Circuit breaker triggered, stopping retries")
			return ErrCircuitBreaker
		}

		err = fn()
		if err == nil {
			log.Printf("Success on attempt %d", i+1)
			cb.Reset() // Reset circuit breaker on success
			return nil
		}

		if i == cfg.MaxRetries-1 {
			log.Println("Max retries reached.")
			return ErrMaxRetries
		}

		cb.RecordFailure()
		backoff := CalculateBackoff(cfg, i)
		log.Printf("Retrying after error: %v (attempt %d/%d) - waiting %v", err, i+1, cfg.MaxRetries, backoff)

		select {
		case <-ctx.Done():
			log.Println("Retry process canceled due to context timeout")
			return ErrContextCanceled
		case <-time.After(backoff):
		}
	}

	log.Println("Max retries reached. Returning error:", err)
	return ErrMaxRetries

}
func RetryMiddleware(cfg Config, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		err := Retry(ctx, func() error {
			rec := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}
			next.ServeHTTP(rec, r)

			if rec.statusCode >= 500 {
				return HTTPError{StatusCode: rec.statusCode, Message: "Retryable HTTP error"}
			}
			return nil
		}, cfg)

		if err != nil {
			if errors.Is(err, ErrCircuitBreaker) {
				log.Println("Circuit breaker triggered for HTTP request")
				http.Error(w, "Circuit breaker triggered - Too many failures", http.StatusTooManyRequests)
			} else if errors.Is(err, ErrMaxRetries) {
				log.Println("Max retries reached for HTTP request")
				http.Error(w, "Service unavailable - Retry limit exceeded", http.StatusServiceUnavailable)
			} else {
				log.Println("Unhandled error in retry middleware:", err)
				http.Error(w, "Unexpected error", http.StatusInternalServerError)
			}
		}

	})
}
func isRetryable(err error, cfg Config) bool {
	if len(cfg.RetryableErrors) > 0 {
		retryableErrs := make(map[error]struct{}, len(cfg.RetryableErrors))
		for _, re := range cfg.RetryableErrors {
			retryableErrs[re] = struct{}{}
		}
		if _, exists := retryableErrs[err]; exists {
			return true
		}
	}

	if httpErr, ok := err.(HTTPError); ok {
		for _, code := range cfg.RetryableStatusCodes {
			if httpErr.StatusCode == code {
				return true
			}
		}
	}
	return false
}

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rr *responseRecorder) WriteHeader(code int) {
	rr.statusCode = code
	rr.ResponseWriter.WriteHeader(code)
}
