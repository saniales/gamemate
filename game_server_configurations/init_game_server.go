package gameServerConfigurations

import (
	"sanino/gamemate/constants"

	"sanino/gamemate/controllers/shared/sockets"

	"github.com/labstack/echo"
)

const (
	TICTACTOE_ID int = 24
)

//InitGameServer creates the game server and binds it to the current structure
//NOTE : this is a prototypical approach. In real implementation the server will act
//       on its own docker container.
func InitGameServer(server *echo.Echo) {
	InitArchives()
	InitCache()

	server.POST(constants.MATCH_WEBSOCKET_CHANNEL_PATH, socketController.HandleChannel)
}
