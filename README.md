
# Go Retry Library

## 📌 Overview

This Go library provides a robust **retry mechanism** with support for:

-   ✅ **Exponential, linear, and fixed backoff strategies**
-   ✅ **Custom retry logic** with user-defined functions
-   ✅ **Circuit breaker protection** to prevent excessive retries
-   ✅ **Retryable HTTP errors** and generic error handling
-   ✅ **Context-aware retries** with timeout and cancellation support

## 🚀 Features

-   **Multiple Backoff Strategies**: Fixed, Linear, Exponential, Custom
-   **Circuit Breaker Support**: Prevents retrying failures indefinitely
-   **Flexible Configuration**: Customize retry behavior easily
-   **Retry Middleware**: Wraps HTTP handlers for automatic retries
-   **Optimized Performance**: Uses maps for fast lookup of retryable errors

----------

## ⚙ Installation

```sh
  go get github.com/devjefster/GoRetry

```

----------

## 🛠 Usage

### 🔹 Basic Example

```go
package main

import (
	"context"
	"errors"
	"fmt"
	"time"
	"github.com/devjefster/GoRetry"
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

```

----------

### 🔹 Using Circuit Breaker

```go
cfg := retry.Config{
	MaxRetries:              5,
	CircuitBreakerThreshold: 2, // Stops retrying after 2 consecutive failures
}

```

----------

### 🔹 Custom Backoff Function

```go
cfg := retry.Config{
	BackoffStrategy: retry.CustomBackoff,
	CustomBackoffFunc: func(attempt int) time.Duration {
		return time.Duration(attempt) * 200 * time.Millisecond
	},
}

```

----------

### 🔹 HTTP Retry Middleware

Wrap your HTTP handlers with retry logic:

```go
mux := http.NewServeMux()
mux.Handle("/data", myHandler)

cfg := retry.DefaultConfig
cfg.RetryableStatusCodes = []int{500, 502, 503, 504}
wrappedHandler := retry.RetryMiddleware(cfg, mux)

http.ListenAndServe(":8080", wrappedHandler)

```

----------

## 🧪 Running Tests

```sh
  go test -v ./...

```

----------

## 📜 License

MIT License

----------

## 🤝 Contributing

1.  Fork the repo
2.  Create a new branch (`git checkout -b feature-name`)
3.  Commit your changes (`git commit -m 'Added new feature'`)
4.  Push to the branch (`git push origin feature-name`)
5.  Open a Pull Request 🎉
