package main

import (
	"GoRetry/retry"
	"log"
	"net/http"
)

func myHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError) // Simulate server failure
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/data", http.HandlerFunc(myHandler))

	cfg := retry.DefaultConfig
	cfg.RetryableStatusCodes = []int{500, 502, 503, 504}

	wrappedHandler := retry.RetryMiddleware(cfg, mux)

	log.Println("Starting server on :8080")
	http.ListenAndServe(":8080", wrappedHandler)
}
