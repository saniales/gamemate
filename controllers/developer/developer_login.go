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
)

//registerDeveloper inserts a developer into the archives.
//
//Returns ID if successfull, error otherwise is filled.
func registerDeveloper(RegTry developerRequests.DevRegistration) (int64, error) {
	authTry := developerRequests.DevAuth{Email: RegTry.Email, Password: RegTry.Password}
	isLoggable, _, err := checkLogin(authTry)
	if err == nil && isLoggable {
		return -1, errors.New("Developer already registered")
	}
	salt := rand.Intn(constants.MAX_NUMBER_SALT)
	saltedPass := controllerSharedFuncs.ConvertToHexString(RegTry.Password + strconv.Itoa(salt))

	stmtQuery, err := configurations.ArchivesPool.Prepare(
		fmt.Sprintf("INSERT INTO developers (developerID, email, hash_pwd, hash_salt) VALUES (NULL, ?, UNHEX(?), ?)"),
	)
	if err != nil {
		return -1, err
	}
	defer stmtQuery.Close()

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
//Additionally returns if found the developerID (V.2)
func checkLogin(AuthTry developerRequests.DevAuth) (bool, int64, error) {

	var num_rows int
	var password_hash string
	var developerID int64
	var salt int

	stmtQuery, err := configurations.ArchivesPool.Prepare("SELECT COUNT(*) AS num_rows, HEX(hash_pwd), hash_salt, developerID FROM developers WHERE email = ? GROUP BY hash_pwd, hash_salt, developerID")
	if err != nil {
		return false, -1, err
	}
	defer stmtQuery.Close()

	result, err := stmtQuery.Query(AuthTry.Email)
	if err != nil {
		return false, -1, err
	}
	if !result.Next() {
		return false, -1, errors.New("Cannot login user")
	}
	err = result.Scan(&num_rows, &password_hash, &salt, &developerID)
	if err != nil {
		return false, -1, err
	}
	salted_pwd := AuthTry.Password + strconv.Itoa(salt)
	salted_hash := controllerSharedFuncs.ConvertToHexString(salted_pwd)
	//fmt.Println("0x" + password_hash)
	//fmt.Println(salted_hash)
	if salted_hash == password_hash {
		return true, developerID, nil
	}
	return false, -1, nil
}

func updateCacheWithSessionDeveloperToken(developerID int64) (string, error) {
	return controllerSharedFuncs.UpdateCacheNewSession(constants.LOGGED_DEVELOPERS_SET, constants.CACHE_REFRESH_INTERVAL, developerID) //can add extra data, see documentation
}

func getDevIDFromSessionToken(token string) (int64, error) {
	return controllerSharedFuncs.GetIDFromSessionSet(constants.LOGGED_DEVELOPERS_SET, token)
}
