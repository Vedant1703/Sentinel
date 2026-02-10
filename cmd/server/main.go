package main

import (
	"encoding/json"
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
	http.Handle("/", middleware.CORSMiddleware(rateLimiter.Middleware(handler)))
	http.Handle("/metrics", middleware.CORSMiddleware(metrics.Handler()))

	// Dynamic Config Endpoint
	http.Handle("/api/config", middleware.CORSMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			Path   string `json:"path"`
			Limit  int    `json:"limit"`
			Window int    `json:"window"` // seconds
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		rateLimiter.UpdateConfig(req.Path, req.Limit, time.Duration(req.Window)*time.Second)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Updated %s: %d reqs / %ds", req.Path, req.Limit, req.Window)
	})))

	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", nil)
}
