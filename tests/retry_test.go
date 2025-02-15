package retry_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"GoRetry/retry"
)

func TestRetryFailure(t *testing.T) {
	fn := func() error {
		return errors.New("permanent failure")
	}

	cfg := retry.DefaultConfig
	cfg.MaxRetries = 3
	ctx := context.Background()

	err := retry.Retry(ctx, fn, cfg)
	if !errors.Is(err, retry.ErrMaxRetries) {
		t.Errorf("Expected max retries error, but got: %v", err)
	}

}

func TestCircuitBreakerStopsRetries(t *testing.T) {
	fn := func() error {
		return errors.New("persistent failure")
	}

	cfg := retry.Config{
		MaxRetries:              5,
		CircuitBreakerThreshold: 2, // Stops retries after 2 failures
	}

	ctx := context.Background()
	err := retry.Retry(ctx, fn, cfg)

	if !errors.Is(err, retry.ErrCircuitBreaker) {
		t.Errorf("Expected circuit breaker error, but got: %v", err)
	}

}

func TestBackoffStrategies(t *testing.T) {
	cfg := retry.Config{
		MaxRetries:     3,
		InitialBackoff: 100 * time.Millisecond,
		MaxBackoff:     1 * time.Second,
	}

	tests := []struct {
		strategy retry.BackoffStrategy
	}{
		{retry.FixedBackoff},
		{retry.LinearBackoff},
		{retry.ExponentialBackoff},
		{retry.CustomBackoff},
	}

	for _, test := range tests {
		cfg.BackoffStrategy = test.strategy
		_ = retry.CalculateBackoff(cfg, 1)
	}
}

func TestRetryMiddleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError) // Simulate server failure
	})

	cfg := retry.DefaultConfig
	cfg.RetryableStatusCodes = []int{500, 502, 503, 504}

	server := httptest.NewServer(retry.RetryMiddleware(cfg, handler))
	defer server.Close()

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusTooManyRequests {
		t.Errorf("Expected 429 Too Many Requests, got: %d", resp.StatusCode)
	}

}
