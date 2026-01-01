package middleware

import (
	"net/http"

	"github.com/Vedant/distributed-rate-limiter/limiter/burst"
)
func RateLimit(limiter *burst.Limiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// Step 1: extract key (IP address)
			key := r.RemoteAddr

			// Step 2: check burst limiter
			if !limiter.Allow(key) {
				w.WriteHeader(http.StatusTooManyRequests)
				w.Write([]byte("Too Many Requests"))
				return
			}

			// Step 3: request allowed → continue
			next.ServeHTTP(w, r)
		})
	}
}
