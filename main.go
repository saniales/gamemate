//Represents the main package of this project and the file to be "executed" at
//first level.
package main

import (
	//serves http requests.

	//echo main package.
	"github.com/labstack/echo/engine/fasthttp" //fast go engine, can be replaced with standard.

	"sanino/gamemate/configurations"
	"sanino/gamemate/constants"
	"sanino/gamemate/controllers"

	_ "sanino/gamemate/libs" //Custom package for this project (structs - API - etc...).
	//Package to interact with Redis DB
)

//Main function of the server : here there are the allowed types of connections
//and their behaviour.
func main() {
	e := configurations.InitServer()
	//Links with redis to permit cache usage.
	configurations.InitCache()
	configurations.InitArchives()
	//defer redisPool.Close()

	e.POST(constants.AUTH_PATH, controllers.HandleAuth)
	e.POST(constants.USER_REGISTRATION_PATH, controllers.HandleRegistration)

	e.Logger().Print(e.Run(fasthttp.New(":8080")))
}
