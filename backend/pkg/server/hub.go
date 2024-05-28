package server

import (
	"fmt"
	"sync"

	"social-network/backend/pkg/db/sqlite"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	connections map[*Client]sqlite.Connection

	// Register requests from the clients.
	connection chan *Client

	// Unregister requests from clients.
	deconnection chan *Client

	mutex sync.Mutex

	previousUser int
}

func NewHub() *Hub {
	return &Hub{
		connections:  make(map[*Client]sqlite.Connection),
		connection:   make(chan *Client),
		deconnection: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.connection:
			h.mutex.Lock()
			if h.previousUser != 0 {
				h.connections[client] = sqlite.Connection{Id: h.previousUser, IsConnected: true}
				h.previousUser = 0
			} else {
				h.connections[client] = sqlite.Connection{Id: 0, IsConnected: false}
			}
			fmt.Println("STATE AFTER USER ENTERED: ", h.connections)
			h.mutex.Unlock()

		case client := <-h.deconnection:
			if client.refreshed {
				h.previousUser = h.connections[client].Id
			}
			close(client.send)
			delete(h.connections, client)
			client.conn.Close()
			fmt.Println("STATE AFTER USER REFRESH OR LEAVE: ", h.connections)
		}
	}
}
