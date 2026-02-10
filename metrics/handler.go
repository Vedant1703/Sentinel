package metrics

import (
	"encoding/json"
	"net/http"
	"sync/atomic"
)

type Response struct {
	Allowed uint64 `json:"allowed_requests"`
	Blocked uint64 `json:"blocked_requests"`
	Errors  uint64 `json:"redis_errors"`
}

func Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		resp := Response{
			Allowed: atomic.LoadUint64(&Allowed), // Use atomic load for safety
			Blocked: atomic.LoadUint64(&Blocked),
			Errors:  atomic.LoadUint64(&Errors),
		}
		json.NewEncoder(w).Encode(resp)
	}
}
