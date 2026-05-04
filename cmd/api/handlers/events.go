package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/degeens/scrobblet/internal/events"
)

func Events(bus *events.Bus) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "streaming not supported", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		ch, unsub := bus.Subscribe()
		defer unsub()

		for {
			select {
			case e := <-ch:
				data, _ := json.Marshal(e)
				fmt.Fprintf(w, "data: %s\n\n", data)
				flusher.Flush()
			case <-r.Context().Done():
				return
			}
		}
	}
}
