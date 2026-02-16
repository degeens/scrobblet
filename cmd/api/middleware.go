package main

import (
	"log/slog"
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
		var (
			ip     = r.RemoteAddr
			proto  = r.Proto
			method = r.Method
			uri    = r.URL.RequestURI()
		)

		slog.Info("Received request", "ip", ip, "proto", proto, "method", method, "uri", uri)

		next.ServeHTTP(w, r)
	})
}
