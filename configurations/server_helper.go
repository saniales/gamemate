package configurations

import (
	"sanino/gamemate/constants"
	"strconv"
	"time"

	"github.com/garyburd/redigo/redis"
)

func collectCacheGarbage() {
	timer := time.NewTicker(constants.CACHE_REFRESH_INTERVAL)
	go func() {
		<-timer.C //waiting for signal from timer channel
		ClearExpiredCache()
	}()
}

//ClearExpiredCache must be run every N minutes with a timer to clear the cache from
//the expired stuff.
func ClearExpiredCache() error {
	conn := CachePool.Get()
	err := removeExpiredLoggedTokens(conn)
	if err != nil {
		return err
	}
	return nil
}

//removeExpiredLoggedTokens removes expired logged_tokens from cache.
func removeExpiredLoggedTokens(conn redis.Conn) error {
	err := conn.Send("ZREMRANGEBYSCORE", "logged_users", "0", strconv.FormatInt(time.Now().UnixNano(), 10))
	if err != nil {
		return err
	}
	return nil
}
