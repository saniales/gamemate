package socketModels

import (
	"errors"
	"math/rand"
	"sanino/gamemate/models/shared/game_server"
	"sanino/gamemate/models/user/data_structures"

	"github.com/gorilla/websocket"
)

//ServerRoom represents a lobby on the server.
type ServerRoom struct {
	ID           int64     //ID of the room.
	hub          SocketHub //Connected users.
	PlayersLeft  int64     //missing players to complete the room.
	MatchStarted bool      //Match started.
}

//NewServerRoom creates a new server room with the specified parameters.
func NewServerRoom(roomID int64, maxPlayers int64) *ServerRoom {
	ret := &ServerRoom{
		ID:           roomID,
		hub:          NewSocketHub(),
		PlayersLeft:  maxPlayers,
		MatchStarted: false,
	}
	return ret
}

//IsFull returns true if it is not possible to add other players,
//false otherwise.
func (receiver *ServerRoom) IsFull() bool {
	return receiver.MatchStarted || receiver.PlayersLeft == 0
}

//AddPlayer Adds a player to the room.
func (receiver *ServerRoom) AddPlayer(p userDataStructs.Player, socket *websocket.Conn) error {
	if receiver.IsFull() {
		return errors.New("Room FULL")
	}
	receiver.hub.Clients[socket] = p
	receiver.PlayersLeft--
	return nil
}

//BroadcastRoomUpdate sends a message to update room of all connected peers.
func (receiver *ServerRoom) BroadcastRoomUpdate(typeOfUpdate string) []error {
	Message := make(map[string]interface{})
	Message["Action"] = typeOfUpdate
	Message["RoomID"] = receiver.ID
	Players := make([]userDataStructs.Player, 0, 0)
	for _, v := range receiver.hub.Clients {
		Players = append(Players, v)
	}
	Message["Players"] = Players
	Message["PlayersLeft"] = receiver.PlayersLeft
	Message["MatchStarted"] = receiver.MatchStarted
	if receiver.MatchStarted {
		conn, firstPlayer := receiver.chooseRandomPlayer()
		Message["FirstPlayer"] = firstPlayer

		if receiver.SetFirstTurn(conn) != nil {
			panic("Cannot set first player")
		}
	}
	return receiver.hub.Broadcast(Message)
}

//chooseRandomPlayer selects a random player from the clients and returns it.
func (receiver *ServerRoom) chooseRandomPlayer() (*websocket.Conn, userDataStructs.Player) {
	keys := make([]*websocket.Conn, len(receiver.hub.Clients))
	i := 0
	for k := range receiver.hub.Clients {
		keys[i] = k
		i++
	}
	key := keys[rand.Intn(i)]
	player, _ := receiver.GetConnectedPlayer(key)
	return key, player
}

//GetConnectedPlayer gets the player connected with the specified socket from this struct.
//If not found second return value is set to false.
func (receiver *ServerRoom) GetConnectedPlayer(conn *websocket.Conn) (userDataStructs.Player, bool) {
	return receiver.hub.GetConnectedPlayer(conn)
}

//IsPlayerTurn returns true if the requesting connection's player is the one who has
//permission to play for this turn, false otherwise.
func (receiver *ServerRoom) IsPlayerTurn(conn *websocket.Conn) bool {
	return receiver.hub.IsPlayerTurn(conn)
}

//NextTurn passes the turn to the nest player.
func (receiver *ServerRoom) NextTurn() {
	receiver.hub.NextTurn()
}

//SetFirstTurn sets, if possible, the first player's turn, otherwise
//it returns an error.
func (receiver *ServerRoom) SetFirstTurn(conn *websocket.Conn) error {
	return receiver.hub.SetFirstTurn(conn)
}

//SendMoveRejected sends to the player which made an invalid move a notification.
func (receiver *ServerRoom) SendMoveRejected(dest *websocket.Conn, CustomData map[string]interface{}) error {
	message := make(map[string]interface{})
	message["Action"] = "MoveRejected"
	message["Cell"] = CustomData["Cell"]
	return dest.WriteJSON(message)
}

//BroadcastNewMove sends a new move to connected players
func (receiver *ServerRoom) BroadcastNewMove(Move map[string]interface{}, result gameServerLogic.Result) []error {
	message := make(map[string]interface{})
	message["Action"] = "NewMove"
	message["NextPlayer"] = receiver.CurrentPlayer()
	message["Cell"] = Move["Cell"]
	message["Symbol"] = Move["Symbol"]
	message["Result"] = result
	return receiver.hub.Broadcast(message)
}

//CurrentPlayer gets the current player which is in turn.
func (receiver *ServerRoom) CurrentPlayer() userDataStructs.Player {
	return receiver.hub.CurrentPlayer()
}
