package socketController

import (
	"errors"
	"net/http"

	"sanino/gamemate/controllers/shared"
	"sanino/gamemate/controllers/user/session_controller"
	"sanino/gamemate/models/shared/responses/errors"
	"sanino/gamemate/models/shared/socket"
	"sanino/gamemate/models/user/data_structures"
	"sanino/gamemate/models/user/requests/out_match"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
)

var (
	currentRoom *socketModels.ServerRoom
	upgrader    = websocket.Upgrader{}
)

//HandleChannel handles a request to create a socket, due to a request.
func HandleChannel(context echo.Context) error {
	ws, err := upgrader.Upgrade(context.Response(), context.Request(), nil)
	if err != nil {
		errorResp := errorResponses.ErrorDetail{}
		context.Logger().Print(err)
		errorResp.FromError(errors.New("Cannot establish connection, retry later"), http.StatusBadRequest)
		return context.JSON(http.StatusBadRequest, &errorResp)
	}
	defer ws.Close()

	for {
		IncomingMessage := make(map[string]string)
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
			val, _ := controllerSharedFuncs.IsValidAPI_Token(IncomingMessage["API_Token"])
			if !val || err != nil {
				errorResp := errorResponses.ErrorDetail{}
				context.Logger().Print(errors.New("Rejected by the system, requestor not valid"))
				errorResp.FromError(err, http.StatusBadRequest)
				return context.JSON(http.StatusBadRequest, &errorResp)
			}
			//check user logged
			userID, _ := sessionController.GetUserIDFromSessionToken(IncomingMessage["SessionToken"])
			if err != nil {
				errorResponse := errorResponses.ErrorDetail{}
				context.Logger().Print(errors.New("Rejected by the system, invalid session"))
				errorResponse.FromError(errors.New("Rejected by the system, invalid session"), http.StatusBadRequest)
				return context.JSON(http.StatusBadRequest, errorResponse)
			}
			//TODO: check if the game is enabled.
			currentRoom = getCurrentRoom()
			if currentRoom.IsFull() {
				errorResponse := errorResponses.ErrorDetail{}
				context.Logger().Print("Room FULL")
				errorResponse.FromError(errors.New("Rejected by the system, the lobby is full"), http.StatusBadRequest)
				return context.JSON(http.StatusBadRequest, errorResponse)
			}

			//TODO: get user detail from cache, for now passing as parameter in request.
			//add ws to the pool
			current := getCurrentRoom()
			current.AddPlayer(userDataStructs.Player{
				ID:       userID,
				Username: request.Username,
			}, ws)
			update := "NewPlayer"
			if current.IsFull() {
				update = "MatchStarted"
			}
			current.BroadcastRoomUpdate(update)
			break
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
