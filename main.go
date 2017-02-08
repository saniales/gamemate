//Represents the main package of this project and the file to be "executed" at
//first level.
package main

import (
	"net/http" //serves http requests.

	"github.com/labstack/echo"                 //echo main package.
	"github.com/labstack/echo/engine/fasthttp" //fast go engine, can be replaced with standard.

	"sanino/gamemate/configurations"
	"sanino/gamemate/constants"
	"sanino/gamemate/controllers"
	"sanino/gamemate/models/request"

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

	e.POST(constants.GET_USER_REQUEST_PATH, func(c echo.Context) error {
		var err error
		user := request.Auth{}
		user.FromForm(c)
		//user.Email = c.FormValue("Email")
		//err = user.InsertIntoDB(redisPool)
		if err != nil {
			return c.JSON(http.StatusCreated, err)
		}
		return c.JSON(http.StatusCreated, user)
	})

	e.Run(fasthttp.New(":8080"))
}
