package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type room struct {
	// forward is a channel that holds incoming messages // that should be forwarded to the other clients.
	forward chan []byte
	// join is a channel for clients wishing to join the room.
	join chan *userClient
	// leave is a channel for clients wishing to leave the room.
	leave chan *userClient
	// clients holds all current clients in this room.
	clients map[*userClient]bool
}

func newRoom() *room {
	return &room{
		forward: make(chan []byte),
		join:    make(chan *userClient),
		leave:   make(chan *userClient),
		clients: make(map[*userClient]bool),
	}
}

func (room *room) run() {
	for {
		select {

		case client := <-room.join:
			room.clients[client] = true
		case client := <-room.leave:
			delete(room.clients, client)
			close(client.send)

		case msg := <-room.forward:
			// forward message to all clients
			for client := range room.clients {
				client.send <- msg
			}
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  socketBufferSize,
	WriteBufferSize: socketBufferSize}

// Adding this method means a room type can now act as http.Handler (implicit signature)
func (room *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	// upgrade the connection to use WS
	socket, err := upgrader.Upgrade(w, req, nil)

	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}

	// create the client
	client := &userClient{
		socket: socket,
		send:   make(chan []byte, messageBufferSize),
		room:   room,
	}

	// join the room
	room.join <- client
	defer func() { room.leave <- client }()

	go client.write() // goroutine
	client.read()     // read to keep alive
}
