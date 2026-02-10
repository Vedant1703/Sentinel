package middleware

import (
	"net"
	"net/http"
)

//	func ExtractKey(r *http.Request) string {
//		if userID := r.Header.Get("X-User-ID"); userID != "" {
//			return "user:" + userID
//		}
//		return "ip:" + r.RemoteAddr
//	}
func ExtractKey(r *http.Request) string {
	if userID := r.Header.Get("X-User-ID"); userID != "" {
		return "user:" + userID
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		ip = r.RemoteAddr
	}

	return "ip:" + ip + ":path:" + r.URL.Path
}
