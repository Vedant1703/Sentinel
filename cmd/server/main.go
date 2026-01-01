package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Vedant/distributed-rate-limiter/limiter/burst"
	"github.com/Vedant/distributed-rate-limiter/middleware"
)

func main() {
	// Create burst limiter
	limiter := burst.NewLimiter(5, 2*time.Second)

	// Demo handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Request allowed")
	})

	// Wrap handler with rate-limiting middleware
	rateLimitedHandler := middleware.RateLimit(limiter)(handler)

	http.Handle("/", rateLimitedHandler)

	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", nil)
}

