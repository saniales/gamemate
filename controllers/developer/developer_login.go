package developerController

import (
	"errors"
	"fmt"
	"math/rand"
	"sanino/gamemate/configurations"
	"sanino/gamemate/constants"
	"sanino/gamemate/controllers/shared"
	"sanino/gamemate/models/developer/requests"
	"strconv"

	"github.com/garyburd/redigo/redis"
)

func registerDeveloper(RegTry developerRequests.DevRegistration) error {
	isLoggable, err := checkLogin(developerRequests.DevAuth{Email: RegTry.Email, Password: RegTry.Password})
	if err == nil {
		return errors.New("Cannot check if user is registered")
	}
	if isLoggable {
		return errors.New("Developer already registered")
	}
	salt := rand.Intn(constants.MAX_NUMBER_SALT)
	saltedPass := RegTry.Password + strconv.Itoa(salt)

	stmtQuery, err := configurations.ArchivesPool.Prepare(
		fmt.Sprintf("INSERT INTO developers (ID, email, password) VALUES (NULL, ?, %s)",
			controllerSharedFuncs.ConvertToHexString(saltedPass)),
	)
	if err != nil {
		return err
	}
	result, err := stmtQuery.Exec(RegTry.Email)
	if err != nil {
		return err
	}
	rowsAff, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAff <= 0 {
		return errors.New("No row affected, possible error with the query")
	}
	return nil
}

//checkLogin searchs for an user with the specified credentials, return error
//some errors occurred with the Archives.
//
//Returns true if found, false otherwise.
func checkLogin(AuthTry developerRequests.DevAuth) (bool, error) {
	panic("OMMERDA")
	//return "", nil
}

func updateCacheWithSessionDeveloperToken(email string) (string, error) {
	return controllerSharedFuncs.UpdateCacheNewSession(constants.LOGGED_DEVELOPERS_SET, email, constants.CACHE_REFRESH_INTERVAL)
}

func getDevEmailFromSessionToken(token string) (string, error) {
	conn := configurations.CachePool.Get()
	email, err := redis.String(conn.Do("GET", "token/"+token+"/"+constants.LOGGED_DEVELOPERS_SET))
	if err != nil {
		return "", err
	}
	if email == "" {
		return "", errors.New("Invalid Session")
	}
	return email, nil
}
