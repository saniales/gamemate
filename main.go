//Represents the main package of this project and the file to be "executed" at
//first level.
package main

import (
	"github.com/labstack/echo/engine/fasthttp" //fast go engine, can be replaced with standard.

	"sanino/gamemate/configurations"
	"sanino/gamemate/constants"
	"sanino/gamemate/controllers/developer"
	"sanino/gamemate/controllers/game_owner"
	"sanino/gamemate/controllers/user/login_controller"
)

//Main function of the server : here there are the allowed types of connections
//and their behaviour.
func main() {
	e := configurations.InitServer()
	//Links with redis to permit cache usage.
	configurations.InitCache()
	configurations.InitArchives()
	//defer redisPool.Close()

	e.POST(constants.AUTH_PATH, loginController.HandleAuth)
	e.POST(constants.USER_REGISTRATION_PATH, loginController.HandleRegistration)

	e.POST(constants.DEVELOPER_AUTH_PATH, developerController.HandleLogin)
	e.POST(constants.DEVELOPER_REGISTRATION_PATH, developerController.HandleRegistration)
	e.POST(constants.DEVELOPER_ADD_API_TOKEN_PATH, developerController.HandleAddAPI_Token)
	e.POST(constants.DEVELOPER_DROP_API_TOKEN, developerController.HandleDropAPI_Token)

	e.POST(constants.GAME_OWNER_AUTH_PATH, gameOwnerController.HandleLogin)
	e.POST(constants.GAME_OWNER_REGISTRATION_PATH, gameOwnerController.HandleRegistration)
	e.POST(constants.GAME_OWNER_ADD_GAME_PATH, gameOwnerController.HandleAddGame)
	e.POST(constants.GAME_OWNER_REMOVE_GAME_PATH, gameOwnerController.HandleRemoveGame)
	e.POST(constants.GAME_ENABLE_DISABLE_PATH, gameOwnerController.HandleGameAction)
	e.POST(constants.GAME_OWNER_GAME_LIST, nil)

	e.Logger().Print(e.Run(fasthttp.New(":8080")))
}
