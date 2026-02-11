package redislimiter

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type Limiter struct {
	client *redis.Client
	script *redis.Script
	limit  int
	window time.Duration
}

func NewLimiter(limit int, window time.Duration) *Limiter {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}

	var opts *redis.Options
	var err error

	if strings.HasPrefix(addr, "redis://") || strings.HasPrefix(addr, "rediss://") {
		opts, err = redis.ParseURL(addr)
		if err != nil {
			log.Fatalf("invalid redis url: %v", err)
		}
	} else {
		opts = &redis.Options{
			Addr: addr,
		}
	}

	rdb := redis.NewClient(opts)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("redis connection failed: %v", err)
	}

	scriptBytes, err := os.ReadFile("redis/ratelimit.lua")
	if err != nil {
		log.Fatalf("failed to load lua script: %v", err)
	}

	return &Limiter{
		client: rdb,
		script: redis.NewScript(string(scriptBytes)),
		limit:  limit,
		window: window,
	}
}

func (l *Limiter) Allow(key string) (bool, error) {
	ctx := context.Background()

	res, err := l.script.Run(
		ctx,
		l.client,
		[]string{key},
		int(l.window.Seconds()),
		l.limit,
	).Int()

	if err != nil {
		return false, err
	}

	return res == 1, nil
}
