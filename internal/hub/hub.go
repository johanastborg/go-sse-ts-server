package hub

import (
	"sync"
)

// Hub manages a pool of clients and broadcasts messages to them.
type Hub[T any] struct {
	mu         sync.RWMutex
	clients    map[chan T]struct{}
	broadcast  chan T
	register   chan chan T
	unregister chan chan T
	done       chan struct{}
}

// NewHub creates a new Hub instance.
func NewHub[T any]() *Hub[T] {
	return &Hub[T]{
		clients:    make(map[chan T]struct{}),
		broadcast:  make(chan T),
		register:   make(chan chan T),
		unregister: make(chan chan T),
		done:       make(chan struct{}),
	}
}

// Run starts the hub's main loop.
func (h *Hub[T]) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = struct{}{}
			h.mu.Unlock()
		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client)
			}
			h.mu.Unlock()
		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client <- message:
				default:
					// If client buffer is full, we might drop messages or handle it differently
				}
			}
			h.mu.RUnlock()
		case <-h.done:
			h.mu.Lock()
			for client := range h.clients {
				close(client)
			}
			h.clients = nil
			h.mu.Unlock()
			return
		}
	}
}

// Broadcast sends a message to all registered clients.
func (h *Hub[T]) Broadcast(msg T) {
	select {
	case h.broadcast <- msg:
	case <-h.done:
	}
}

// Register adds a new client channel to the hub.
func (h *Hub[T]) Register() chan T {
	ch := make(chan T, 100)
	h.register <- ch
	return ch
}

// Unregister removes a client channel from the hub.
func (h *Hub[T]) Unregister(ch chan T) {
	h.unregister <- ch
}

// Stop shuts down the hub.
func (h *Hub[T]) Stop() {
	close(h.done)
}
