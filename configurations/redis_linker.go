package configurations

import (
	"github.com/garyburd/redigo/redis" //Package to interact with Redis DB
)

const (
	REDIS_HOST            string = "127.0.0.1" //Address of the Redis server
	REDIS_PORT            string = "6379"      //Port of the Redis server
	REDIS_MAX_CONNECTIONS int    = 12000       //Max number of simultaneous connections allowed to the Redis server
)

//RedisPool represents the pool which connects to the cache (using Redis).
var RedisPool *redis.Pool

//CacheInitialized is true if the Pool has been initialized at least one time.
var CacheInitialized = false

//InitCache creates a Redis communication point for the API cache.
func InitCache() {
	RedisPool = redis.NewPool(func() (redis.Conn, error) {
		c, err := redis.Dial("tcp", REDIS_HOST+":"+string(REDIS_PORT))
		if err != nil {
			return nil, err
		}
		return c, err
	}, REDIS_MAX_CONNECTIONS)
	CacheInitialized = true
}
