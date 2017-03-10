package developerController

import (
	"errors"
	"fmt"

	"sanino/gamemate/configurations"
	"sanino/gamemate/constants"
	"sanino/gamemate/controllers/shared"

	"github.com/garyburd/redigo/redis"
)

//getAPITokensOfDeveloper gets tokens from the archives or cache.
func getAPITokensOfDeveloper(developerID int64) ([]string, bool, error) {
	tokens, err := checkAPITokenListInCache(developerID)
	if err != nil {
		tokens, err = getAPITokenListFromArchives(developerID)
		if err != nil {
			if updateCacheWithTokenList(developerID, tokens) != nil {
				return tokens, false, nil
			}
		}
		return tokens, true, nil
	}
	return tokens, true, nil
}

func checkAPITokenListInCache(developerID int64) ([]string, error) {
	conn := configurations.CachePool.Get()
	defer conn.Close()
	key := fmt.Sprintf(constants.DEVELOPER_TOKEN_LIST_IN_CACHE, developerID)
	_, err := conn.Do("EXPIRE", key, constants.CACHE_REFRESH_INTERVAL.Seconds())
	if err != nil {
		return nil, err
	}
	tokens, err := redis.Strings(conn.Do("SMEMBERS", key))
	return tokens, err
}

func updateCacheWithTokenList(developerID int64, tokens []string) error {
	conn := configurations.CachePool.Get()
	defer conn.Close()
	key := fmt.Sprintf(constants.DEVELOPER_TOKEN_LIST_IN_CACHE, developerID)
	err := conn.Send("MULTI")
	if err != nil {
		return err
	}
	err = conn.Send("EXPIRE", key, constants.CACHE_REFRESH_INTERVAL.Seconds())
	if err != nil {
		conn.Do("DISCARD")
		return err
	}
	for _, val := range tokens {
		err = conn.Send("SADD", key, val)
		if err != nil {
			conn.Do("DISCARD")
			return err
		}
	}
	_, err = conn.Do("EXEC")
	if err != nil {
		conn.Do("DISCARD")
		return err
	}
	return nil
}

func addTokenToCacheList(developerID int64, token string) error {
	conn := configurations.CachePool.Get()
	defer conn.Close()
	key := fmt.Sprintf(constants.DEVELOPER_TOKEN_LIST_IN_CACHE, developerID)
	_, err := conn.Do("SADD", key, token)
	return err
}

func getAPITokenListFromArchives(developerID int64) ([]string, error) {
	stmtQuery, err := configurations.ArchivesPool.Prepare("SELECT token FROM API_Tokens WHERE developerID = ?")
	if err != nil {
		return nil, errors.New("Cannot get tokens, query prepare error => " + err.Error())
	}
	defer stmtQuery.Close()
	rows, err := stmtQuery.Query(developerID)
	if err != nil {
		return nil, errors.New("Cannot get tokens, query params error => " + err.Error())
	}

	tokens := make([]string, 10)

	for rows.Next() {
		var token string
		if rows.Err() != nil {
			return nil, errors.New("Cannot get tokens, query row error => " + err.Error())
		}
		if rows.Scan(token) != nil {
			return nil, errors.New("Cannot get tokens, query row-scan error => " + err.Error())
		}
		tokens = append(tokens, token)
	}

	return tokens, nil
}

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

	result, err := redis.Int64(conn.Do("SISMEMBER", constants.API_TOKENS_SET, token))
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

	//the cache is valid for 24 hours, if an app is not used it should not be in cache.
	err := conn.Send("SADD", constants.API_TOKENS_SET, token)
	if err != nil {
		return err
	}

	err = conn.Flush()
	if err != nil {
		return err
	}
	return nil
}

//removeAPI_TokenFromCache removes from the Cache the specified API_Token.
//
//Return error if did not manage to update the cache.
func removeAPI_TokenFromCache(token string) error {
	conn := configurations.CachePool.Get()
	defer conn.Close()
	err := conn.Send("SREM", constants.API_TOKENS_SET, token)
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
	stmtQuery, err := configurations.ArchivesPool.Prepare(
		"SELECT COUNT(token) FROM API_Tokens WHERE token = ? AND enabled = 1"
	)
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
	fmt.Print(num_rows)
	if num_rows > 0 {
		err = updateCacheWithAPI_Token(token)
		if err != nil { //did not update cache but the request has been satisfied.
			return true, nil
		}
		return true, nil
	}
	return false, errors.New("API Token not found")
}

//addAPI_TokenInArchives adds a token linked to the specified developer to the archives.
func addAPI_TokenInArchives(developerID int64) (string, error) {
	token := controllerSharedFuncs.GenerateToken()
	//TODO: find a way to handle duplicates. or leave the query fail and retry.
	stmtQuery, err := configurations.ArchivesPool.Prepare(
		fmt.Sprintf("INSERT INTO API_Tokens (developerID, token, enabled) VALUES (?, %s, 1)",
			controllerSharedFuncs.ConvertToHexString(token)),
	)
	if err != nil {
		return "", err
	}
	defer stmtQuery.Close()
	result, err := stmtQuery.Exec(developerID)
	if err != nil {
		return "", err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return "", err
	}
	if rows <= 0 {
		return "", errors.New("No Row Affected, possible problem with the query")
	}
	return token, nil
}

//removeAPI_TokenFromArchives removes a token from the Archives.
//
//Request is valid only if the API Token to remove is owned by the requestor.
func removeAPI_TokenFromArchives(developerID int64, token string) error {
	stmtQuery, err := configurations.ArchivesPool.Prepare(
		fmt.Sprintf("UPDATE API_Tokens SET enabled = 0 WHERE token = %s AND developerID = ?",
			controllerSharedFuncs.ConvertToHexString(token)),
	)
	if err != nil {
		return err
	}
	defer stmtQuery.Close()

	result, err := stmtQuery.Exec(developerID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows <= 0 {
		return errors.New("No Row Affected, possible problem with the query or developerID is not the owner")
	}
	return nil
}
