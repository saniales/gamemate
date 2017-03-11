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

//GenerateToken creates a valid session token, without "0x" HEX marker.
func GenerateToken() string {
	return strings.Replace(ConvertToHexString(strconv.FormatInt(time.Now().UnixNano(), 10)), "0x", "", 1)
}

//UpdateCacheNewSession Updates the cache SessionSet with the email and the expiration.
//
//ID represents the ID of the entity to put
//values represents additional data to create hashmaps,
//they must be in the form ["key1", "value1", "key2", "value2", ..., and so on];
//can be empty.
func UpdateCacheNewSession(SessionSet string, expiration time.Duration, ID int64, values ...string) (string, error) {
	var token string
	conn := configurations.CachePool.Get()
	defer conn.Close()
	tokenOk := false
	for i := 0; i < 10; i++ {
		token = GenerateToken()
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

	err = conn.Send("HMSET", SessionSet+"/with_token/"+token, "ID", ID, values)
	//if set when a user logons with an expired key it is removed from cache and set
	if err != nil {
		return constants.INVALID_TOKEN, err
	}

	err = conn.Send("EXPIRE", SessionSet+"/with_token/"+token, expiration.Seconds())
	//if set when a user logons with an expired key it is removed from cache and set
	if err != nil {
		return constants.INVALID_TOKEN, err
	}

	command := fmt.Sprintf("%s/%d", SessionSet, ID)
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
