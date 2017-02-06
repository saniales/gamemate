package configurations

import (
	"github.com/labstack/echo" //echo main package.
	"github.com/labstack/echo/middleware"
)

//Configures the server for fresh start with the default configuration.
func InitServer() *echo.Echo {
	server := echo.New()
	// Middleware
	server.Use(middleware.Logger())
	server.Use(middleware.Recover())
	//CORS
	server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.POST},
	}))
	return server
}
