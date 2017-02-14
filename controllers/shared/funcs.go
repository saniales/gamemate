package controllerSharedFuncs

import (
	"crypto/sha512"
	"encoding/hex"
	"errors"
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
func UpdateCacheNewSession(SessionSet string, email string, expiration time.Duration) (string, error) {
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
	err = conn.Send("SET", SessionSet+"/"+email+"/token", token, "EX", int(expiration.Seconds()))
	//if set when a user logons with an expired key it is removed from cache and set
	if err != nil {
		return constants.INVALID_TOKEN, err
	}

	err = conn.Send("SET", "token/"+token+"/"+SessionSet, email, "EX", int(expiration.Seconds()))
	//if set when a user logons with an expired key it is removed from cache and set
	if err != nil {
		return constants.INVALID_TOKEN, err
	}

	err = conn.Send("EXEC")
	if err != nil {
		return constants.INVALID_TOKEN, err
	}

	err = conn.Flush()
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
