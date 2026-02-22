package middleware

import (
	"log/slog"
	"net"
	"net/http"
	"net/url"
)

func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		proto := r.Proto
		method := r.Method
		sanitizedURI := sanitizeURI(*r.URL)

		ip := r.RemoteAddr
		if host, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
			ip = host
		}

		slog.Info("Received request", "proto", proto, "method", method, "uri", sanitizedURI, "ip", ip)

		next.ServeHTTP(w, r)
	})
}

func sanitizeURI(u url.URL) string {
	q := u.Query()

	// Hide sensitive OAuth query parameters if present
	if q.Has("code") {
		q.Set("code", "hidden")
	}
	if q.Has("state") {
		q.Set("state", "hidden")
	}

	u.RawQuery = q.Encode()
	return u.RequestURI()
}
