package main

import (
	"log/slog"
	"net"
	"net/http"

	"golang.org/x/time/rate"
)

func rateLimit(r int, b int) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		// Allow r requests per second sustained, up to b in a burst
		limiter := rate.NewLimiter(rate.Limit(r), b)

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !limiter.Allow() {
				http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method := r.Method
		uri := r.URL.RequestURI()
		proto := r.Proto

		ip := r.RemoteAddr
		if host, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
			ip = host
		}

		slog.Info("Received request", "method", method, "uri", uri, "proto", proto, "ip", ip)

		next.ServeHTTP(w, r)
	})
}
