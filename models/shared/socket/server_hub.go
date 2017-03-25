package socketModels

import (
	"github.com/gorilla/websocket"

	"sanino/gamemate/models/user/data_structures"
)

//SocketHub represents a hub with multiple connections to clients.
type SocketHub struct {
	Clients map[*websocket.Conn]userDataStructs.Player //Represents the pool of connected clients.
}

//NewSocketHub creates a new Hub for a room.
func NewSocketHub() SocketHub {
	return SocketHub{
		Clients: make(map[*websocket.Conn]userDataStructs.Player),
	}
}

//AddClient adds a connection to the hub, linking it with the player.
func (receiver *SocketHub) AddClient(client *websocket.Conn, player userDataStructs.Player) {
	receiver.Clients[client] = player
}

//Broadcast sends to all connected peers a message
func (receiver *SocketHub) Broadcast(message interface{}) []error {
	errors := make([]error, 0, 1)
	for conn := range receiver.Clients {
		err := conn.WriteJSON(message)
		if err != nil {
			errors = append(errors, err)
		}
	}
	if len(errors) == 0 {
		return nil
	}
	return errors
}
