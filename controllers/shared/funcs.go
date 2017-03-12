package controllerSharedFuncs

import (
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"sanino/gamemate/configurations"
	"sanino/gamemate/constants"
	"strconv"
	"strings"
	"time"

	"github.com/garyburd/redigo/redis"
)

//generateToken creates a valid session token.
func GenerateToken() string {
	return ConvertToHexString(strconv.FormatInt(time.Now().UnixNano(), 10))
}

//UpdateCacheNewSession Updates the cache SessionSet with the email and the expiration.
//
//ID represents the ID of the entity to put
//values represents additional data to create hashmaps,
//they must be in the form ["key1", "value1", "key2", "value2", ..., and so on];
//can be empty.
func UpdateCacheNewSession(SessionSet string, expiration time.Duration, ID int64, values ...interface{}) (string, error) {
	var token string
	conn := configurations.CachePool.Get()
	defer conn.Close()
	tokenOk := false
	for i := 0; i < 10; i++ {
		token = strings.Replace(GenerateToken(), "0x", "", 1)
		_, err := redis.String(conn.Do("ZSCORE", SessionSet, token))
		if err != nil {
			if err.Error() != "(redigo: nil returned)" { //if nil is ok, break
				tokenOk = true
				break
			}
			return constants.INVALID_TOKEN, err
		}
	}
	if !tokenOk { //tried a lot of times to find a free session token
		return constants.INVALID_TOKEN, errors.New("Rejected by the system, retry")
	}
	err := conn.Send("MULTI")
	if err != nil {
		return constants.INVALID_TOKEN, err
	}

	err = conn.Send("ZADD", SessionSet, time.Now().Add(expiration).UnixNano(), token) //sets the cache for the token expire 30 minutes
	if err != nil {
		return constants.INVALID_TOKEN, err
	}

	var args []interface{}
	args = append(args, SessionSet+"/with_token/"+token, "ID", ID)
	args = append(args, values...)
	err = conn.Send("HMSET", args...)
	//if set when a user logons with an expired key it is removed from cache and set
	if err != nil {
		return constants.INVALID_TOKEN, err
	}

	err = conn.Send("EXPIRE", SessionSet+"/with_token/"+token, expiration.Seconds())
	if err != nil {
		return constants.INVALID_TOKEN, err
	}

	command := fmt.Sprintf("%s/with_id/%d", SessionSet, ID)
	err = conn.Send("SET", command, token, "EX", int(expiration.Seconds()))
	//if set when a user logons with an expired key it is removed from cache and set
	if err != nil {
		return constants.INVALID_TOKEN, err
	}

	_, err = conn.Do("EXEC")
	if err != nil {
		return constants.INVALID_TOKEN, err
	}

	return token, nil
}

//ConvertToHexString converts a string to a SHA512 Representation string.
func ConvertToHexString(source string) string {
	hash := sha512.New()
	io.WriteString(hash, source)
	return "0x" + strings.ToUpper(hex.EncodeToString(hash.Sum(nil)))
}

//GetIDFromSessionSet gets a generic ID of an entity from its set in cache and session token.
func GetIDFromSessionSet(SessionSet string, Token string) (int64, error) {
	command := fmt.Sprintf("%s/with_token/%s", SessionSet, Token)

	conn := configurations.CachePool.Get()
	ID, err := redis.Int64(conn.Do("HGET", command, "ID"))
	if err != nil {
		return -1, fmt.Errorf("Invalid Session : command = %s, response = %d, error = %v", command, ID, err)
	}
	if ID == 0 {
		return -1, fmt.Errorf("Invalid Session : command = %s, response = %d, error = %v", command, ID, err)
	}
	return ID, nil
}
