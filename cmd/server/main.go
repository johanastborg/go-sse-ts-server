package main

import (
	"log"
	"net/http"

	"github.com/johanastborg/go-sse-ts-server/internal/feed"
	"github.com/johanastborg/go-sse-ts-server/internal/hub"
	"github.com/johanastborg/go-sse-ts-server/internal/sse"
)

func main() {
	// 1. Initialize the Fan-out Hub
	h := hub.NewHub[feed.DataPoint]()
	go h.Run()
	defer h.Stop()

	// 2. Initialize the Sine Wave Feed (50Hz)
	// Frequency = 1Hz, Amplitude = 10, Noise = 0.5
	producer := feed.NewSineProducer(1.0, 10.0, 0.5)
	dataCh := producer.Start()
	defer producer.Stop()

	// 3. Pipe feed to hub
	go func() {
		for dp := range dataCh {
			h.Broadcast(dp)
		}
	}()

	// 4. Setup SSE Handler
	sseHandler := sse.NewHandler(h)

	// 5. Setup HTTP Server
	mux := http.NewServeMux()
	mux.Handle("/stream", sseHandler)

	// Health check or index
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Go SSE Time Series Server Running... Connect to /stream"))
	})

	log.Println("Server starting on :8080...")
	log.Println("Connect with: curl -N http://localhost:8080/stream")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
