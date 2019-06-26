package main

import (
	"net/http"

	uuid "github.com/satori/go.uuid"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

func wsPage(res http.ResponseWriter, req *http.Request) {
	conn, err := upgrader.Upgrade(res, req, nil)
	if err != nil {
		http.NotFound(res, req)
		return
	}
	clientID, _ := uuid.NewV4()
	client := &Client{id: clientID.String(), socket: conn, send: make(chan []byte)}

	manager.register <- client

	go client.read()
	go client.write()
}
