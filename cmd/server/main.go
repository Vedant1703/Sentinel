package main

import (
	"fmt"
	"time"

	"github.com/Vedant/distributed-rate-limiter/limiter/burst"
)

func main() {
	limiter := burst.NewLimiter(3, 2*time.Second)

	key := "user-1"

	for i := 1; i <= 5; i++ {
		allowed := limiter.Allow(key)
		fmt.Println("Request", i, "allowed:", allowed)
	}
}
