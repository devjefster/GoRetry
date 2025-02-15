package main

import (
	"GoRetry/retry"
	"context"
	"errors"
	"fmt"
	"time"
)

func unstableCustom() error {
	if time.Now().UnixNano()%3 != 0 {
		return errors.New("random failure")
	}
	return nil
}

func main() {
	cfg := retry.Config{
		BackoffStrategy: retry.CustomBackoff,
		CustomBackoffFunc: func(attempt int) time.Duration {
			return time.Duration(attempt) * 200 * time.Millisecond
		},
		MaxRetries: 5,
	}

	ctx := context.Background()
	err := retry.Retry(ctx, unstableCustom, cfg)
	if err != nil {
		fmt.Println("Operation failed:", err)
	} else {
		fmt.Println("Operation succeeded")
	}
}
