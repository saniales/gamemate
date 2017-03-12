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
	//      So for debugging purposes using HTTP bacause cannot use 443 (8080) assigned me from committant.
	//server.Pre(middleware.HTTPSRedirect())
	server.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `{"time":"${time_rfc3339_nano}","remote_ip":"${remote_ip}","host":"${host}",` +
			`"method":"${method}", "uri":"${uri}","status":${status}, "latency":${latency},` +
			`"latency_human":"${latency_human}","bytes_in":${bytes_in},` +
			`"bytes_out":${bytes_out}}` + "\n",
	}))
	server.Use(middleware.Recover())
	server.Logger.SetLevel(log.INFO)
	//CORS
	server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.POST},
	}))
	collectCacheGarbage(server)
	return server
}
