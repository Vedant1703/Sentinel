package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Vedant/distributed-rate-limiter/metrics"

	"github.com/Vedant/distributed-rate-limiter/config"
	burstlimiter "github.com/Vedant/distributed-rate-limiter/limiter/burst"
	"github.com/Vedant/distributed-rate-limiter/middleware"
)

func main() {

	// Burst limiter (local)
	burstLimiter := burstlimiter.NewLimiter(20, 50*time.Millisecond)

	// Configurable per-route rules
	cfg := config.Config{
		Routes: map[string]config.Rule{
			"/login":  {Limit: 5, Window: time.Minute},
			"/search": {Limit: 50, Window: time.Minute},
		},
		Default: config.Rule{
			Limit:  100,
			Window: time.Minute,
		},
	}

	rateLimiter := middleware.NewRateLimitMiddleware(
		burstLimiter,
		cfg,
		true, // fail-open
	)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Request allowed")
	})
	http.Handle("/", rateLimiter.Middleware(handler))
	http.Handle("/metrics", metrics.Handler())

	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", nil)
}
