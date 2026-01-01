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
func (l *Limiter) Allow(key string) (bool, error) {
	return l.allow(key), nil
}

func (l *Limiter) allow(key string) bool {	
	now := time.Now()
	l.mu.Lock()
	defer l.mu.Unlock()

	timestamps := l.requests[key]
	valid := make([]time.Time, 0, len(timestamps))

	for _, ts := range timestamps {
		if now.Sub(ts) <= l.window {
			valid = append(valid, ts)
		}
	}

	if len(valid) >= l.limit {
		l.requests[key] = valid
		return false
	}

	valid = append(valid, now)
	l.requests[key] = valid
	return true
}
