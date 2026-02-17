package main

import (
	"log/slog"
	"net"
	"net/http"
)

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
