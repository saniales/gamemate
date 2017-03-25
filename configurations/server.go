package configurations

import (
	"math/rand"
	"time"

	"github.com/labstack/echo" //echo main package.
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
)

//InitServer configures the server for fresh start with the default configuration.
func InitServer() *echo.Echo {
	rand.Seed(time.Now().UTC().UnixNano())
	server := echo.New()
	//Cache TLS Certificates
	//server.AutoTLSManager.Cache = autocert.DirCache("/tmp/gamemate/.cache")

	// Middleware

	//NOTE : HTTPS is valid only on port 443 for ACME generator, have to generate it manually
	//      So for debugging purposes using HTTPS self signed certificate manually bacause
	//      cannot use 443 (8080) assigned me from committant.
	//server.Pre(middleware.HTTPSRedirect())
	server.Use(middleware.Secure())
	server.Use(middleware.Logger())
	server.Use(middleware.Recover())
	server.Use(middleware.Gzip())
	server.Use(middleware.BodyLimit("1M"))

	server.Logger.SetLevel(log.INFO)
	//CORS
	server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.POST, echo.GET},
	}))
	collectCacheGarbage(server)
	return server
}
