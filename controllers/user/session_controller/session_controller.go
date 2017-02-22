package sessionController

import (
	"errors"
	"sanino/gamemate/configurations"

	"github.com/garyburd/redigo/redis"
)

//GetUserIDFromSessionToken gets the user ID from the session token in cache.
func GetUserIDFromSessionToken(token string) (int64, error) {
	conn := configurations.CachePool.Get()
	ID, err := redis.Int64(conn.Do("HMGET", "users/with_token/"+token+"/", "ID"))
	if err != nil {
		return -1, err
	}
	if ID == 0 {
		return -1, errors.New("Invalid Session")
	}
	return ID, nil
}
