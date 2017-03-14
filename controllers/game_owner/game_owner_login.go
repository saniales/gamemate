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
)

func registerOwner(RegTry gameOwnerRequests.GameOwnerRegistration) (int64, error) {
	AuthTry := gameOwnerRequests.GameOwnerAuth{Email: RegTry.Email, Password: RegTry.Password}
	isLoggable, _, err := checkLogin(AuthTry)
	if err != nil {
		return -1, errors.New("Login Error Check => " + err.Error())
	}
	if isLoggable {
		return -1, errors.New("Owner already registered, cannot register again")
	}
	salt := rand.Intn(constants.MAX_NUMBER_SALT)
	saltedPass := controllerSharedFuncs.ConvertToHexString(RegTry.Password + strconv.Itoa(salt))

	stmtQuery, err := configurations.ArchivesPool.Prepare(
		fmt.Sprintf("INSERT INTO game_owners (ownerID, email, hash_pwd, hash_salt) VALUES (NULL, ?, UNHEX(?), ?)"),
	)
	if err != nil {
		return -1, err
	}
	result, err := stmtQuery.Exec(RegTry.Email, saltedPass, salt)
	if err != nil {
		return -1, err
	}
	rowsAff, err := result.RowsAffected()
	if err != nil {
		return -1, err
	}
	if rowsAff <= 0 {
		return -1, errors.New("No row affected, possible error with the query")
	}

	insertID, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}
	return insertID, nil
}

//checkLogin searchs for an user with the specified credentials, return error
//some errors occurred with the Archives.
//
//Returns true if found, false otherwise.
func checkLogin(AuthTry gameOwnerRequests.GameOwnerAuth) (bool, int64, error) {
	var num_rows int
	var password_hash string
	var ownerID int64
	var salt int

	stmtQuery, err := configurations.ArchivesPool.Prepare("SELECT COUNT(*) AS num_rows, HEX(hash_pwd), hash_salt, ownerID FROM game_owners WHERE email = ? GROUP BY hash_pwd, hash_salt, ownerID")
	if err != nil {
		return false, -1, err
	}
	defer stmtQuery.Close()

	result, err := stmtQuery.Query(AuthTry.Email)
	if err != nil {
		return false, -1, err
	}
	if !result.Next() { //user not found, so no error
		return false, -1, nil
	}
	err = result.Scan(&num_rows, &password_hash, &salt, &ownerID)
	if err != nil {
		return false, -1, err
	}
	salted_pwd := AuthTry.Password + strconv.Itoa(salt)
	salted_hash := controllerSharedFuncs.ConvertToHexString(salted_pwd)
	//fmt.Println("0x" + password_hash)
	//fmt.Println(salted_hash)
	if salted_hash == password_hash {
		return true, ownerID, nil
	}
	return false, -1, nil
}

func updateCacheWithSessionOwnerToken(ownerID int64) (string, error) {
	return controllerSharedFuncs.UpdateCacheNewSession(constants.LOGGED_OWNERS_SET, constants.CACHE_REFRESH_INTERVAL, ownerID)
}

func getOwnerIDFromSessionToken(token string) (int64, error) {
	return controllerSharedFuncs.GetIDFromSessionSet(constants.LOGGED_OWNERS_SET, token)
}
