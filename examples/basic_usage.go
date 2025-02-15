package main

import (
	"GoRetry/retry"
	"context"
	"errors"
	"fmt"
	"time"
)

func unstableOperation() error {
	if time.Now().UnixNano()%2 == 0 {
		return errors.New("temporary failure")
	}
	return nil
}

func main() {
	cfg := retry.DefaultConfig
	cfg.MaxRetries = 3

	ctx := context.Background()
	err := retry.Retry(ctx, unstableOperation, cfg)
	if err != nil {
		fmt.Println("Operation failed:", err)
	} else {
		fmt.Println("Operation succeeded")
	}
}
