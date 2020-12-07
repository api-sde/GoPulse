package main

import "github.com/gorilla/websocket"

/// Represent a single chatting user
type userClient struct {
	socket *websocket.Conn

	// channel on which messages are sent
	send chan []byte

	// room is where the user is chatting
	room *room
}

func (cli *userClient) read() {
	defer cli.socket.Close()

	for {
		_, msg, err := cli.socket.ReadMessage()

		if err != nil {
			return
		}
		cli.room.forward <- msg

	}
}

func (cli *userClient) write() {
	defer cli.socket.Close()

	for msg := range cli.send {
		err := cli.socket.WriteMessage(websocket.TextMessage, msg)

		if err != nil {
			return
		}
	}

}
