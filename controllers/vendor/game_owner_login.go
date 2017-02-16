package gameOwnerController

import (
	"errors"
	"fmt"
	"math/rand"
	"sanino/gamemate/configurations"
	"sanino/gamemate/constants"
	"sanino/gamemate/controllers/shared"
	"sanino/gamemate/models/game_owner/requests"
	"strconv"

	"github.com/garyburd/redigo/redis"
)

func registerOwner(RegTry gameOwnerRequests.GameOwnerRegistration) error {
	isLoggable, err := checkLogin(gameOwnerRequests.GameOwnerAuth{Email: RegTry.Email, Password: RegTry.Password})
	if err == nil {
		return errors.New("Cannot check if user is registered")
	}
	if isLoggable {
		return errors.New("Developer already registered")
	}
	salt := rand.Intn(constants.MAX_NUMBER_SALT)
	saltedPass := RegTry.Password + strconv.Itoa(salt)

	stmtQuery, err := configurations.ArchivesPool.Prepare(
		fmt.Sprintf("INSERT INTO game_owners (ID, email, password) VALUES (NULL, ?, %s)",
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
func checkLogin(AuthTry gameOwnerRequests.GameOwnerAuth) (bool, error) {

	var num_rows int
	var password_hash string
	var salt int

	stmtQuery, err := configurations.ArchivesPool.Prepare("SELECT COUNT(*) AS num_rows, HEX(hash_pwd), hash_salt FROM developers WHERE email = ? GROUP BY hash_pwd, hash_salt")
	if err != nil {
		return false, err
	}
	defer stmtQuery.Close()

	result, err := stmtQuery.Query(AuthTry.Email)
	if err != nil {
		return false, err
	}
	if !result.Next() {
		return false, errors.New("Cannot login user, Database error")
	}
	err = result.Scan(&num_rows, &password_hash, &salt)
	if err != nil {
		return false, err
	}
	salted_pwd := AuthTry.Password + strconv.Itoa(salt)
	salted_hash := controllerSharedFuncs.ConvertToHexString(salted_pwd)
	//fmt.Println("0x" + password_hash)
	//fmt.Println(salted_hash)
	return salted_hash == "0x"+password_hash, nil
}

func updateCacheWithSessionOwnerToken(email string) (string, error) {
	return controllerSharedFuncs.UpdateCacheNewSession(constants.LOGGED_OWNERS_SET, email, constants.CACHE_REFRESH_INTERVAL)
}

func getOwnerIDFromSessionToken(token string) (int64, error) {
	conn := configurations.CachePool.Get()
	ID, err := redis.Int64(conn.Do("GET", "token/"+token+"/"+constants.LOGGED_OWNERS_SET))
	if err != nil {
		return "", err
	}
	if ID == 0 {
		return -1, errors.New("Invalid Session")
	}
	return ID, nil
}
