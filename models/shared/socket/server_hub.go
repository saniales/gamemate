package socketModels

import (
	"errors"

	"github.com/gorilla/websocket"

	"sanino/gamemate/models/user/data_structures"

	log "github.com/labstack/gommon/log"
)

//SocketHub represents a hub with multiple connections to clients.
type SocketHub struct {
	Clients          map[*websocket.Conn]userDataStructs.Player //Represents the pool of connected clients.
	Turns            []*websocket.Conn                          //Represents the keys used to get the turns during a match.
	CurrentTurnIndex int                                        //Represents the index in turns array to tell which player must make the move.
}

//NewSocketHub creates a new Hub for a room.
func NewSocketHub() SocketHub {
	return SocketHub{
		Clients:          make(map[*websocket.Conn]userDataStructs.Player),
		Turns:            make([]*websocket.Conn, 0),
		CurrentTurnIndex: -1,
	}
}

//AddClient adds a connection to the hub, linking it with the player.
func (receiver *SocketHub) AddClient(client *websocket.Conn, player userDataStructs.Player) {
	receiver.Clients[client] = player
	receiver.Turns = append(receiver.Turns, client)
}

//Broadcast sends to all connected peers a message.
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

//GetConnectedPlayer gets the player connected with the specified socket from this struct.
//If not found second return value is set to false.
func (receiver *SocketHub) GetConnectedPlayer(conn *websocket.Conn) (userDataStructs.Player, bool) {
	val, ok := receiver.Clients[conn]
	return val, ok
}

//IsPlayerTurn returns true if the requesting connection's player is the one who has
//permission to play for this turn, false otherwise.
func (receiver *SocketHub) IsPlayerTurn(conn *websocket.Conn) bool {
	log.Debug(receiver.CurrentTurnIndex)
	return receiver.Clients[conn] == receiver.Clients[receiver.Turns[receiver.CurrentTurnIndex]]
}

//NextTurn passes the turn to the nest player.
func (receiver *SocketHub) NextTurn() {
	receiver.CurrentTurnIndex = (receiver.CurrentTurnIndex + 1) % len(receiver.Turns)
}

//SetFirstTurn sets, if possible, the first player's turn, otherwise
//it returns an error.
func (receiver *SocketHub) SetFirstTurn(conn *websocket.Conn) error {
	for i, v := range receiver.Turns {
		if v == conn {
			receiver.CurrentTurnIndex = i
			return nil
		}
	}
	return errors.New("Player not found")
}

//CurrentPlayer gets the current player which is in turn.
func (receiver *SocketHub) CurrentPlayer() userDataStructs.Player {
	return receiver.Clients[receiver.Turns[receiver.CurrentTurnIndex]]
}
