package metrics

import (
	"fmt"
	"net/http"
)

func Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "allowed_requests %d\n", Allowed)
		fmt.Fprintf(w, "blocked_requests %d\n", Blocked)
		fmt.Fprintf(w, "redis_errors %d\n", Errors)
	}
}
