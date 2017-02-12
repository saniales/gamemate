package developerController

import (
	"errors"
	"fmt"

	"sanino/gamemate/configurations"

	"github.com/garyburd/redigo/redis"
)

//IsValidAPI_Token Provides a control for forged requests with fake API_Tokens
//
//Returns true if the token is valid, false otherwise.
func IsValidAPI_Token(token string) (bool, error) {
	var msgCache, msgArchives string
	isInCache, errCache := checkAPI_TokenInCache(token)
	if !isInCache {
		if errCache != nil {
			msgCache = errCache.Error()
		} else {
			msgCache = "No Error"
		}
		isInArchives, errArchives := checkAPI_TokenInArchives(token)
		if !isInArchives {
			if errArchives != nil {
				msgArchives = errArchives.Error()
			} else {
				msgArchives = "No Error"
			}
			return false, errors.New("Check API Error: \"" +
				msgCache + "\" Message from Cache and \"" +
				msgArchives + "\" Message from Archives.")
		}
		//TODO: OK, but cache has errors (which should have been fixed, but this feature is not implemented yet)
		return true, errors.New("Check API Error: Cache says \"" + msgCache + "\"")
	}
	return true, nil
}

//checkAPI_TokenInCache searchs for an API_Token in the cache.
//
//Return true if found, false otherwise.
func checkAPI_TokenInCache(token string) (bool, error) {
	conn := configurations.CachePool.Get()
	defer conn.Close()

	result, err := redis.Int64(conn.Do("SISMEMBER", "API_Tokens", token))
	if err != nil {
		return false, err
	}

	return result == 1, nil
}

//updateCacheWithAPI_Token updates the Cache with the specified API_Token.
//
//Return error if did not manage to update the cache.
func updateCacheWithAPI_Token(token string) error {
	conn := configurations.CachePool.Get()
	defer conn.Close()
	err := conn.Send("SADD", "API_Tokens", token)
	if err != nil {
		return err
	}
	err = conn.Flush()
	if err != nil {
		return err
	}
	return nil
}

//checkAPI_TokenInArchives checks for the existance of the token in the archives
//and, if found, updates the cache.
//
//Return true if found, false otherwise.
func checkAPI_TokenInArchives(token string) (bool, error) {
	stmtQuery, err := configurations.ArchivesPool.Prepare("SELECT COUNT(token) FROM API_Tokens WHERE token = ?")
	if err != nil {
		return false, err
	}
	defer stmtQuery.Close()
	result, err := stmtQuery.Query(token)
	if err != nil {
		return false, err
	}
	if !result.Next() {
		return false, errors.New("Check API error (archives) : Empty Table, Query with errors (should report 0 when item is not in table)")
	}

	var num_rows int64
	result.Scan(&num_rows)

	if num_rows > 0 {
		err = updateCacheWithAPI_Token(token)
		if err != nil {
			return true, fmt.Errorf("API_Token %s has been found in archives, however, Cache has not been updated", token)
		}
		return true, nil
	}
	return false, nil
}
