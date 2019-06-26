package main

import (
	"encoding/json"
)

// ClientManager represents a client manager
type ClientManager struct {
	broadcast  chan []byte
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
}

func (manager *ClientManager) send(message []byte, ignore *Client) {
	for conn := range manager.clients {
		if conn != ignore {
			conn.send <- message
		}
	}
}

func (manager *ClientManager) start() {
	for {
		select {
		case conn := <-manager.register:
			manager.clients[conn] = true
			jsonMessage, _ := json.Marshal(&Message{Content: "/A new socket has connected."})
			manager.send(jsonMessage, conn)

		case conn := <-manager.unregister:
			if _, connected := manager.clients[conn]; connected == true {
				close(conn.send)
				delete(manager.clients, conn)
				jsonMessage, _ := json.Marshal(&Message{Content: "/A socket has been disconnected."})
				manager.send(jsonMessage, conn)
			}

		case message := <-manager.broadcast:
			for conn := range manager.clients {
				select {
				case conn.send <- message:
				default:
					close(conn.send)
					delete(manager.clients, conn)
				}
			}
		}
	}
}
