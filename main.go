//Represents the main package of this project and the file to be "executed" at
//first level.
package main

import (
    "net/http" //serves http requests.
    "github.com/labstack/echo" //echo main package.
    "github.com/labstack/echo/engine/fasthttp" //fast go engine, can be replaced.
    //"strconv" //To convert numbers from Strings and viceversa.
    . "sanino/gamemate/libs" //Custom package for this project (structs - API - etc...).
    . "sanino/gamemate/constants" //Custom package for project constants (e.g. PATHS to connect to API)
    "sanino/gamemate/configurations"
    "github.com/garyburd/redigo/redis" //Package to interact with Redis DB
)

//Main function of the server : here there are the allowed types of connections
//and their behaviour.
func main() {
    var e *echo.Echo = configurations.InitServer()
    //Links with redis to permit cache usage.
    redisPool := configurations.LinkRedis()
    //defer redisPool.Close()

    e.POST(AUTH_PATH, controllers.HandleAuth)

    e.POST(GET_USER_REQUEST_PATH, func(c echo.Context) error {
        var err error = nil
        user := new(Auth)
        user.Username = c.FormValue("Username")
        //user.Email = c.FormValue("Email")
        //err = user.InsertIntoDB(redisPool)
        if err != nil {
          return c.JSON(http.StatusCreated, err)
        } else {
          return c.JSON(http.StatusCreated, user)
        }
    })
    e.GET(ROOT_PATH, func(c echo.Context) error {
        conn := redisPool.Get()
        defer conn.Close()
        lista_utenti, _ := redis.StringMap(conn.Do("HGETALL", "user:1"))
        return c.JSON(http.StatusCreated, lista_utenti)
    })
    e.Run(fasthttp.New(":8080"))
}
