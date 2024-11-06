package sse

import (
	"sync"

	"hafiedh.com/downloader/internal/pkg/constants"
)

type (
	SSEClient struct {
		Message chan constants.DefaultResponse
		ID      string
	}

	SSEHub struct {
		clients      map[string]*SSEClient
		addClient    chan *SSEClient
		removeClient chan *SSEClient
		mu           sync.RWMutex
	}
)

func NewSSEHub() *SSEHub {
	return &SSEHub{
		clients:      make(map[string]*SSEClient),
		addClient:    make(chan *SSEClient),
		removeClient: make(chan *SSEClient),
	}
}

func (h *SSEHub) Run() {
	for {
		select {
		case client := <-h.addClient:
			h.mu.Lock()
			h.clients[client.ID] = client
			h.mu.Unlock()
		case client := <-h.removeClient:
			h.mu.Lock()
			delete(h.clients, client.ID)
			close(client.Message)
			h.mu.Unlock()
		}
	}
}

func (h *SSEHub) AddClient(client *SSEClient) {
	h.mu.Lock()
	h.addClient <- client
	h.mu.Unlock()
}

func (h *SSEHub) RemoveClient(client *SSEClient) {
	h.mu.Lock()
	h.removeClient <- client
	h.mu.Unlock()
}

func (h *SSEHub) BroadcastMessage(message constants.DefaultResponse) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, client := range h.clients {
		client.Message <- message
	}
}
