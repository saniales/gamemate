package socketController

import (
	"errors"
	"net/http"

	"sanino/gamemate/controllers/shared"
	"sanino/gamemate/controllers/user/session_controller"
	"sanino/gamemate/models/shared/game_server"
	"sanino/gamemate/models/shared/responses/errors"
	"sanino/gamemate/models/shared/socket"
	"sanino/gamemate/models/user/data_structures"
	"sanino/gamemate/models/user/requests/out_match"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
)

var (
	currentRoom    *socketModels.ServerRoom
	currentChecker gameServerLogic.MoveChecker
	upgrader       = websocket.Upgrader{}
)

//HandleChannel handles a request to create a socket, due to a request.
func HandleChannel(context echo.Context) error {
	ws, err := upgrader.Upgrade(context.Response(), context.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	for {
		IncomingMessage := make(map[string]interface{})
		//read json message
		err = ws.ReadJSON(&IncomingMessage)
		if err != nil {
			return err
		}
		switch IncomingMessage["Type"] {
		case "GetRoom":
			request := outMatchRequests.GetRoom{}
			err = request.FromMap(IncomingMessage)
			if err != nil {
				return err
			}
			//check api token
			val, _ := controllerSharedFuncs.IsValidAPI_Token(IncomingMessage["API_Token"].(string))
			if err != nil {
				return err
			}
			if !val {
				return errors.New("Invalid API Token")
			}
			//check user logged
			userID, err := sessionController.GetUserIDFromSessionToken(IncomingMessage["SessionToken"].(string))
			if err != nil {
				return err
			}
			//TODO: check if the game is enabled.
			currentRoom = getCurrentRoom()

			//TODO: get user detail from cache, for now passing as parameter in request.
			//add ws to the pool
			current := getCurrentRoom()
			playerToAdd := userDataStructs.Player{
				ID:       userID,
				Username: request.Username,
			}
			err = current.AddPlayer(playerToAdd, ws)
			if err != nil {
				errorResponse := errorResponses.ErrorDetail{}
				context.Logger().Print(err)
				errorResponse.FromError(errors.New("Rejected by the system, the lobby is full"), http.StatusBadRequest)
				return context.JSON(http.StatusBadRequest, errorResponse)
			}
			update := "RoomUpdate"
			if current.IsFull() {
				current.MatchStarted = true
				currentChecker := getCurrentChecker()
			}
			current.BroadcastRoomUpdate(update)
			break
		case "Move":
			//check api token
			val, _ := controllerSharedFuncs.IsValidAPI_Token(IncomingMessage["API_Token"].(string))
			if err != nil {
				return err
			}
			if !val {
				return errors.New("Invalid API Token")
			}
			//check user logged
			_, err := sessionController.GetUserIDFromSessionToken(IncomingMessage["SessionToken"].(string))
			if err != nil {
				return err
			}
			currentRoom = getCurrentRoom()
			if currentRoom.IsPlayerTurn(ws) {
				var CustomData map[string]interface{} = IncomingMessage["CustomData"].(map[string]interface{})
				currentChecker := getCurrentChecker()
				if currentChecker.IsValidMove(CustomData) {

				}
			}
			//get player from conn
			//if player's turn
			//  checkmove
			//else moverejected
			//checkmove:
			//  const {x, y, symbol, player} = request
			//  if grid[x][y] != EMPTY_CELL
			//    Rejected
			//  else if !playersTurn
			//    Rejected
			//  return moveOK and broadcast.
		case "":
		default:
			return errors.New("No Type Defined")
		}
	}
}

//getCurrentRoom() prototypical implementation : returns the current open room.
//NOTE: in final implementation gathers it from a pool of rooms.
func getCurrentRoom() *socketModels.ServerRoom {
	if currentRoom == nil {
		currentRoom = socketModels.NewServerRoom(1, 2)
	}
	return currentRoom
}

func getCurrentChecker() gameServerLogic.MoveChecker {
	if currentChecker == nil {
		currentChecker = gameServerLogic.NewTicTacToeChecker()
	}
	return currentChecker
}
