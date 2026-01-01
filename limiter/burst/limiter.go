package burst

import (
	"sync"
	"time"
)

type Limiter struct {
	mu       sync.Mutex
	requests map[string][]time.Time
	limit    int
	window   time.Duration
}

func NewLimiter(limit int, window time.Duration) *Limiter {
	return &Limiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
}

func (l *Limiter) Allow(key string) bool {
	now := time.Now()

	l.mu.Lock()
	defer l.mu.Unlock()

	// Get existing timestamps
	timestamps := l.requests[key]

	// Keep only timestamps inside the window
	valid := timestamps[:0]
	for _, ts := range timestamps {
		if now.Sub(ts) <= l.window {
			valid = append(valid, ts)
		}
	}

	// Check limit
	if len(valid) >= l.limit {
		l.requests[key] = valid
		return false
	}

	// Allow request
	valid = append(valid, now)
	l.requests[key] = valid
	return true
}

