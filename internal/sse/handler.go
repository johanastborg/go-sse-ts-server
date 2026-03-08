package sse

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/johanastborg/go-sse-ts-server/internal/feed"
	"github.com/johanastborg/go-sse-ts-server/internal/hub"
)

// Handler handles SSE connections.
type Handler struct {
	hub *hub.Hub[feed.DataPoint]
}

// NewHandler creates a new SSE handler.
func NewHandler(h *hub.Hub[feed.DataPoint]) *Handler {
	return &Handler{hub: h}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// SSE Headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Flush the headers immediately
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	// Register client with hub
	clientCh := h.hub.Register()
	defer h.hub.Unregister(clientCh)

	// Send initial message
	fmt.Fprintf(w, ": connected\n\n")
	flusher.Flush()

	// Listening for data or disconnect
	for {
		select {
		case data, ok := <-clientCh:
			if !ok {
				return
			}
			jsonData, err := json.Marshal(data)
			if err != nil {
				continue
			}
			fmt.Fprintf(w, "data: %s\n\n", jsonData)
			flusher.Flush()
		case <-r.Context().Done():
			return
		}
	}
}
