package main

import (
	"fmt"
	"net/http"
	"time"

	burstlimiter "github.com/Vedant/distributed-rate-limiter/limiter/burst"
	redislimiter "github.com/Vedant/distributed-rate-limiter/limiter/redis"
	"github.com/Vedant/distributed-rate-limiter/middleware"
)

func main() {
	burstLimiter := burstlimiter.NewLimiter(
		20,
		50*time.Millisecond,
	)

	redisLimiter := redislimiter.NewLimiter(
		10,
		60*time.Second,
	)

	rateLimitMiddleware := middleware.NewRateLimitMiddleware(
		burstLimiter,
		redisLimiter,
		true,
	)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Request allowed")
	})

	http.Handle("/", rateLimitMiddleware.Middleware(handler))

	fmt.Println("Server running on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
