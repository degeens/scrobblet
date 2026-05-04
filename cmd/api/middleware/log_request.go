package middleware

import (
	"log/slog"
	"net"
	"net/http"
	"net/url"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func (r *statusRecorder) Flush() {
	if f, ok := r.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

func (r *statusRecorder) Unwrap() http.ResponseWriter {
	return r.ResponseWriter
}

func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		proto := r.Proto
		method := r.Method
		sanitizedURI := sanitizeURI(*r.URL)

		ip := r.RemoteAddr
		if host, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
			ip = host
		}

		statusRecorder := &statusRecorder{ResponseWriter: w}

		next.ServeHTTP(statusRecorder, r)

		slog.Info("Received request", "proto", proto, "method", method, "uri", sanitizedURI, "status", statusRecorder.status, "ip", ip)
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
