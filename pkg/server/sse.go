package server

import (
	"fmt"
	"net/http"
	"time"
)

// handleSSE handles streaming server-sent events to frontend clients.
func (s *Server) handleSSE(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported by client", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Subscribe to event bus
	ch, unsubscribe := s.bus.Subscribe()
	defer unsubscribe()

	// Send initial connection handshake
	fmt.Fprintf(w, "event: connected\ndata: {\"status\":\"connected\",\"timestamp\":%d}\n\n", time.Now().UnixMilli())
	flusher.Flush()

	// Keepalive ticker to prevent intermediate proxy timeout
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-r.Context().Done():
			return
		case <-ticker.C:
			fmt.Fprintf(w, ": heartbeat\n\n")
			flusher.Flush()
		case msg, ok := <-ch:
			if !ok {
				return
			}
			fmt.Fprintf(w, "event: %s\ndata: %s\n\n", msg.Event, msg.Data)
			flusher.Flush()
		}
	}
}
