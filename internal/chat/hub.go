package chat

import (
	"context"
	"sync"

	"github.com/gorilla/websocket"
)

type Hub struct {
	Connections map[string]map[*websocket.Conn]bool
	Consumers   map[string]context.CancelFunc
	mutex       sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		Connections: make(map[string]map[*websocket.Conn]bool),
		Consumers:   make(map[string]context.CancelFunc),
	}
}

func (h *Hub) AddConnection(challengeID string, conn *websocket.Conn, startConsumer func()) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if h.Connections[challengeID] == nil {
		h.Connections[challengeID] = make(map[*websocket.Conn]bool)
	}

	h.Connections[challengeID][conn] = true

	if len(h.Connections[challengeID]) == 1 {
		startConsumer()
	}
}

func (h *Hub) RemoveConnection(challengeID string, conn *websocket.Conn) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	delete(h.Connections[challengeID], conn)

	if len(h.Connections[challengeID]) == 0 {
		if cancel, ok := h.Consumers[challengeID]; ok {
			cancel()
			delete(h.Consumers, challengeID)
		}
	}
}

func (h *Hub) Broadcast(challengeID string, message []byte) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	for conn := range h.Connections[challengeID] {
		conn.WriteMessage(websocket.TextMessage, message)
	}
}
