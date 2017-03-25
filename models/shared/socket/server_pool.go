package socketModels

import (
	"sanino/gamemate/models/user/data_structures"

	"github.com/gorilla/websocket"
)

type ServerRoom struct {
	ID           int64     //ID of the room.
	hub          SocketHub //Connected users.
	PlayersLeft  int64     //missing players to complete the room.
	MatchStarted bool      //Match started.
}

func (receiver *ServerRoom) IsFull() bool {
	return receiver.MatchStarted || receiver.PlayersLeft == 0
}

func NewServerRoom(roomID int64, maxPlayers int64) *ServerRoom {
	ret := &ServerRoom{
		ID:           roomID,
		hub:          SocketHub{},
		PlayersLeft:  maxPlayers - 1,
		MatchStarted: false,
	}
	return ret
}

//AddPlayer Adds a player to the room.
func (receiver *ServerRoom) AddPlayer(p userDataStructs.Player, socket *websocket.Conn) {
	receiver.hub.Clients[socket] = p
}

func (receiver *ServerRoom) BroadCastMatchStarted() []error {
	receiver.MatchStarted = true
	Message := make(map[string]interface{})
	Message["RoomID"] = receiver.ID
	Players := make([]userDataStructs.Player, 0, 0)
	for _, v := range receiver.hub.Clients {
		Players = append(Players, v)
	}
	Message["Players"] = Players
	return receiver.hub.BroadCast(Message)
}
