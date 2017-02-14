package configurations

import (
	"math/rand"
	"sanino/gamemate/constants"
	"time"

	"github.com/labstack/echo" //echo main package.
	"github.com/labstack/echo/middleware"
)

//InitServer configures the server for fresh start with the default configuration.
func InitServer() *echo.Echo {
	rand.Seed(time.Now().UTC().UnixNano())
	server := echo.New()
	// Middleware
	server.Use(middleware.Logger())
	server.Use(middleware.Recover())
	//CORS
	server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.POST},
	}))
	server.SetDebug(constants.DEBUG)
	collectCacheGarbage()
	return server
}
