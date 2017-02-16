//Represents the main package of this project and the file to be "executed" at
//first level.
package main

import (
	"github.com/labstack/echo/engine/fasthttp" //fast go engine, can be replaced with standard.

	"sanino/gamemate/configurations"
	"sanino/gamemate/constants"
	"sanino/gamemate/controllers/developer"
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

	e.POST(constants.VENDOR_AUTH_PATH, nil)
	e.POST(constants.VENDOR_REGISTRATION_PATH, nil)
	e.POST(constants.VENDOR_ADD_GAME_PATH, nil)
	e.POST(constants.VENDOR_REMOVE_GAME_PATH, nil)
	e.POST(constants.VENDOR_GAME_LIST, nil)

	e.Logger().Print(e.Run(fasthttp.New(":8080")))
}
