package configurations

import (
	"sanino/gamemate/constants"
	"strconv"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/labstack/echo"
)

func collectCacheGarbage(server *echo.Echo) {
	timer := time.NewTicker(constants.CACHE_REFRESH_INTERVAL)
	go func() {
		<-timer.C //waiting for signal from timer channel
		err := ClearExpiredCache()
		if err != nil {
			server.Logger().Print("Error during garbage collection of the cache => " + err.Error())
		}
	}()
	go ExpireDaylyCacheAtMidnight(server)
}

//ClearExpiredCache must be run every N minutes with a timer to clear the cache from
//the expired stuff.
func ClearExpiredCache() error {
	conn := CachePool.Get()
	defer conn.Close()
	return removeExpiredLoggedTokens(conn)
}

//removeExpiredLoggedTokens removes expired logged_tokens from cache.
func removeExpiredLoggedTokens(conn redis.Conn) error {
	err := conn.Send("ZREMRANGEBYSCORE", constants.LOGGED_OWNERS_SET, "0", strconv.FormatInt(time.Now().UnixNano(), 10))
	if err != nil {
		return err
	}

	err = conn.Send("ZREMRANGEBYSCORE", constants.LOGGED_DEVELOPERS_SET, "0", strconv.FormatInt(time.Now().UnixNano(), 10))
	if err != nil {
		return err
	}

	err = conn.Send("ZREMRANGEBYSCORE", constants.LOGGED_USERS_SET, "0", strconv.FormatInt(time.Now().UnixNano(), 10))
	if err != nil {
		return err
	}

	return nil
}

//ExpireDaylyCacheAtMidnight expires all cache items which must be reset at midnight (UTC).
//
//This is a routine.
func ExpireDaylyCacheAtMidnight(server *echo.Echo) {
	date := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.UTC).Add(time.Hour * 24)
	timer := time.NewTimer(date.Sub(time.Now()))
	for {
		<-timer.C
		err := removeAPITokens()
		if err != nil {
			server.Logger().Print("Error during garbage collection, cannot remove API Tokens")
		}
		timer.Stop()
		date = time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.UTC).Add(time.Hour * 24)
		timer = time.NewTimer(date.Sub(time.Now()))
	}
}

func removeAPITokens() error {
	conn := CachePool.Get()
	defer conn.Close()
	return conn.Send("DEL", constants.API_TOKENS_SET)
}
