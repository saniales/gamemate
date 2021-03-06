//Represents the main package of this project and the file to be "executed" at
//first level.
package main

import (
	"net/http"

	"sanino/gamemate/configurations"
	"sanino/gamemate/constants"
	"sanino/gamemate/controllers/developer"
	"sanino/gamemate/controllers/game_owner"
	"sanino/gamemate/controllers/user/login_controller"
	"sanino/gamemate/controllers/user/out_game_controller"
	"sanino/gamemate/game_server_configurations"

	"github.com/labstack/echo"
)

type dummy struct {
	Message string `json:"message" xml:"message" form:"message"`
}

//Main function of the server : here there are the allowed messages of connections
//and their behaviour.
func main() {
	e := configurations.InitServer()
	//Links with redis to permit cache usage.
	configurations.InitCache()
	configurations.InitArchives()
	//defer redisPool.Close()
	e.POST("/", func(c echo.Context) error { return c.JSON(http.StatusOK, &dummy{Message: "Up and running..."}) })

	e.POST(constants.AUTH_PATH, loginController.HandleAuth)
	e.POST(constants.USER_REGISTRATION_PATH, loginController.HandleRegistration)

	e.POST(constants.USER_ALL_GAME_LIST_PATH, outGameController.HandleAllGamesForUser)
	e.POST(constants.USER_ENABLED_GAME_LIST_PATH, outGameController.HandleMyEnabledGamesForUser)

	e.POST(constants.DEVELOPER_AUTH_PATH, developerController.HandleLogin)
	e.POST(constants.DEVELOPER_REGISTRATION_PATH, developerController.HandleRegistration)
	e.POST(constants.DEVELOPER_TOKEN_LIST_PATH, developerController.HandleAllTokensForDeveloper)
	e.POST(constants.DEVELOPER_ADD_API_TOKEN_PATH, developerController.HandleAddAPI_Token)
	e.POST(constants.DEVELOPER_DROP_API_TOKEN_PATH, developerController.HandleDropAPI_Token)

	e.POST(constants.GAME_OWNER_AUTH_PATH, gameOwnerController.HandleLogin)
	e.POST(constants.GAME_OWNER_REGISTRATION_PATH, gameOwnerController.HandleRegistration)
	e.POST(constants.GAME_OWNER_ADD_GAME_PATH, gameOwnerController.HandleAddGame)
	e.POST(constants.GAME_OWNER_REMOVE_GAME_PATH, gameOwnerController.HandleRemoveGame)
	e.POST(constants.GAME_ENABLE_DISABLE_PATH, gameOwnerController.HandleGameAction)
	e.POST(constants.GAME_OWNER_GAME_LIST_PATH, gameOwnerController.HandleShowMyGames)

	gameServerConfigurations.InitGameServer(e)

	e.Logger.Print(e.Start(":8080"))
	//e.Logger.Print(e.StartTLS(":8080", "/opt/GameMate/certs/cert.pem", "/opt/GameMate/certs/key.pem"))
}
