package socket

import (
	"github.com/googollee/go-socket.io"
	"log"
)

// list of socket.io clients
var Clients = make(map[string]socketio.Socket)

func Broadcast(client string, msgType string, msg string) {
	err := Clients[client].Emit(msgType, msg)
	if err != nil {
		log.Println("Error occurred during sending message through websocket:", err)
	}
}
