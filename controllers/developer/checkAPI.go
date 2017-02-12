package developerController

import (
	"sanino/gamemate/configurations"

	"github.com/garyburd/redigo/redis"
)

//IsValidAPI_Token Provides a control for forged requests with fake API_Tokens
func IsValidAPI_Token(token string) (bool, error) {
	conn := configurations.CachePool.Get()
	defer conn.Close()

	result, err := redis.Int64(conn.Do("SISMEMBER", "API_Tokens", token))
	if err != nil {
		return false, err
	}
	return result == 1, nil
}
