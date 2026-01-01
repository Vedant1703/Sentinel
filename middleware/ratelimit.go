package middleware

import (
	"net"
	"net/http"

	burstlimiter "github.com/Vedant/distributed-rate-limiter/limiter/burst"
	redislimiter "github.com/Vedant/distributed-rate-limiter/limiter/redis"
)

type RateLimitMiddleware struct {
	burstLimiter *burstlimiter.Limiter
	redisLimiter *redislimiter.Limiter
	failOpen     bool
}

func NewRateLimitMiddleware(
	burst *burstlimiter.Limiter,
	redis *redislimiter.Limiter,
	failOpen bool,
) *RateLimitMiddleware {
	return &RateLimitMiddleware{
		burstLimiter: burst,
		redisLimiter: redis,
		failOpen:     failOpen,
	}
}

func (rl *RateLimitMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ip = r.RemoteAddr
		}

		key := "ip:" + ip

		if !rl.burstLimiter.Allow(key) {
			http.Error(w, "rate limit exceeded (burst)", http.StatusTooManyRequests)
			return
		}

		allowed, err := rl.redisLimiter.Allow(key)
		if err != nil {
			if !rl.failOpen {
				http.Error(w, "rate limit unavailable", http.StatusServiceUnavailable)
				return
			}
		} else if !allowed {
			http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
