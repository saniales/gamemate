package configurations

import (
	"math/rand"
	"time"

	"github.com/labstack/echo" //echo main package.
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"golang.org/x/crypto/acme/autocert"
)

//InitServer configures the server for fresh start with the default configuration.
func InitServer() *echo.Echo {
	rand.Seed(time.Now().UTC().UnixNano())
	server := echo.New()
	server.Logger.SetLevel(log.INFO)
	//Cache TLS Certificates
	server.AutoTLSManager.Cache = autocert.DirCache("/tmp/gamemate/.cache")

	// Middleware
	server.Pre(middleware.HTTPSRedirect())
	server.Use(middleware.Logger())
	server.Use(middleware.Recover())
	//CORS
	server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.POST},
	}))
	collectCacheGarbage(server)
	return server
}
