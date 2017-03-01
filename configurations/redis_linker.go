package configurations

import (
	"fmt"

	"github.com/garyburd/redigo/redis" //Package to interact with Redis DB
)

const (
	REDIS_HOST            string = "127.0.0.1" //Address of the Redis server
	REDIS_PORT            int    = 6379        //Port of the Redis server
	REDIS_MAX_CONNECTIONS int    = 12000       //Max number of simultaneous connections allowed to the Redis server
)

//CachePool represents the pool which connects to the cache (using Redis).
var CachePool *redis.Pool

//CacheInitialized is true if the Pool has been initialized at least one time.
var CacheInitialized = false

//InitCache creates a Redis communication point for the API cache.
func InitCache() {
	CachePool = redis.NewPool(func() (redis.Conn, error) {
		c, err := redis.Dial("tcp", fmt.Sprintf("%s:%d", REDIS_HOST, REDIS_PORT))
		if err != nil {
			c = nil
		}
		return c, err
	}, REDIS_MAX_CONNECTIONS)
	CacheInitialized = true
}
