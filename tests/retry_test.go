package retry_test

import (
	"GoRetry/retry"
	"context"
	"errors"
	"testing"
)

func TestRetrySuccess(t *testing.T) {
	attempts := 0
	fn := func() error {
		attempts++
		if attempts < 3 {
			return errors.New("temporary failure")
		}
		return nil
	}

	cfg := retry.DefaultConfig
	cfg.MaxRetries = 5
	cfg.CircuitBreakerThreshold = 5 // Ensure the breaker doesn't stop retries too early
	ctx := context.Background()

	err := retry.Retry(ctx, fn, cfg)
	if err != nil {
		t.Errorf("Expected success, but got error: %v", err)
	}
}

func TestRetryFailure(t *testing.T) {
	fn := func() error {
		return errors.New("permanent failure")
	}

	cfg := retry.DefaultConfig
	cfg.MaxRetries = 3
	ctx := context.Background()

	err := retry.Retry(ctx, fn, cfg)
	if err == nil {
		t.Errorf("Expected failure, but got success")
	}
}

func TestCircuitBreakerResetsAfterSuccess(t *testing.T) {
	attempts := 0
	fn := func() error {
		attempts++
		if attempts < 3 {
			return errors.New("failure")
		}
		return nil
	}

	cfg := retry.Config{
		MaxRetries:              5,
		CircuitBreakerThreshold: 2,
		SuccessReset:            2, // Requires 2 consecutive successes to reset
	}

	ctx := context.Background()
	err := retry.Retry(ctx, fn, cfg)

	if err != nil {
		t.Errorf("Expected success after circuit breaker reset, got: %v", err)
	}
}
