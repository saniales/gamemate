package developerController

import (
	"errors"
	"fmt"
	"strings"

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
		tokens = append(tokens, strings.Replace(token, "0x", "", 1))
	}

	return tokens, nil
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

//removeAPI_TokenFromCache removes from the Cache the specified API_Token.
//
//Return error if did not manage to update the cache.
func removeAPI_TokenFromCache(token string) error {
	conn := configurations.CachePool.Get()
	defer conn.Close()
	_, err := conn.Do("SREM", constants.API_TOKENS_SET, token)
	if err != nil {
		return err
	}
	return nil
}
