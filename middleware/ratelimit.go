package middleware

import (
	"log"
	"net/http"

	"github.com/Vedant/distributed-rate-limiter/config"
	"github.com/Vedant/distributed-rate-limiter/limiter"
	burstlimiter "github.com/Vedant/distributed-rate-limiter/limiter/burst"
	redislimiter "github.com/Vedant/distributed-rate-limiter/limiter/redis"
	"github.com/Vedant/distributed-rate-limiter/metrics"
)

type RateLimitMiddleware struct {
	burstLimiter  limiter.RateLimiter
	redisLimiters map[string]limiter.RateLimiter
	cfg           config.Config
	failOpen      bool
}

func NewRateLimitMiddleware(
	burst *burstlimiter.Limiter,
	cfg config.Config,
	failOpen bool,
) *RateLimitMiddleware {

	redisLimiters := make(map[string]limiter.RateLimiter)

	for route, rule := range cfg.Routes {
		redisLimiters[route] = redislimiter.NewLimiter(rule.Limit, rule.Window)
	}

	// default limiter
	redisLimiters["default"] = redislimiter.NewLimiter(
		cfg.Default.Limit,
		cfg.Default.Window,
	)

	return &RateLimitMiddleware{
		burstLimiter:  burst,
		redisLimiters: redisLimiters,
		cfg:           cfg,
		failOpen:      failOpen,
	}
}

func (rl *RateLimitMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		key := ExtractKey(r)
		path := r.URL.Path

		// 1️⃣ Burst limiter
		allowedBurst, err := rl.burstLimiter.Allow(key)
		if err != nil {
			log.Printf("BURST ERROR key=%s path=%s err=%v\n", key, path, err)
			metrics.IncErrors()
			if !rl.failOpen {
				http.Error(w, "rate limiter unavailable", http.StatusServiceUnavailable)
				return
			}
		} else if !allowedBurst {
			log.Printf("BLOCKED burst key=%s path=%s\n", key, path)
			metrics.IncBlocked()
			http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		// 2️⃣ Choose route limiter
		limiter, ok := rl.redisLimiters[path]
		if !ok {
			limiter = rl.redisLimiters["default"]
		}

		allowed, err := limiter.Allow(key)
		if err != nil {
			log.Printf("REDIS ERROR key=%s path=%s err=%v\n", key, path, err)
			metrics.IncErrors()

			if !rl.failOpen {
				http.Error(w, "rate limiter unavailable", http.StatusServiceUnavailable)
				return
			}
		} else if !allowed {
			log.Printf("BLOCKED redis key=%s path=%s\n", key, path)
			metrics.IncBlocked()
			http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		metrics.IncAllowed()
		next.ServeHTTP(w, r)
	})
}
