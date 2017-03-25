package socketModels

import (
	"github.com/gorilla/websocket"

	"sanino/gamemate/models/user/data_structures"
)

type SocketHub struct {
	Clients map[*websocket.Conn]userDataStructs.Player //Represents the pool of connected clients.
}

func NewSocketHub() SocketHub {
	return SocketHub{
		Clients: make(map[*websocket.Conn]userDataStructs.Player),
	}
}

func (receiver *SocketHub) AddClient(client *websocket.Conn, player userDataStructs.Player) {
	receiver.Clients[client] = player
}

func (receiver *SocketHub) BroadCast(message interface{}) []error {
	errors := make([]error, 0, 1)
	for conn, _ := range receiver.Clients {
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
