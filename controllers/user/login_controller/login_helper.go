package loginController

import (
	"errors"
	"math/rand"

	"sanino/gamemate/configurations"
	"sanino/gamemate/constants"
	"sanino/gamemate/controllers/shared"
	"sanino/gamemate/models/user/requests/login"

	"strconv"
)

//insertIntoArchives (without check of previous insertions, only error reporting)
//inserts a new User into the archives, doing the salty & hashy work.
func insertIntoArchives(RegTry loginRequests.Registration) error {
	salt := rand.Intn(constants.MAX_NUMBER_SALT)
	saltedPass := RegTry.Password + strconv.Itoa(salt)

	stmtQuery, err := configurations.ArchivesPool.Prepare("INSERT INTO users (id, username, hash_pwd, hash_salt, email, birthday, gender) VALUES (NULL, ?," + controllerSharedFuncs.ConvertToHexString(saltedPass) + ", ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmtQuery.Close()

	result, err := stmtQuery.Exec(RegTry.Username, salt, RegTry.Email, RegTry.Birthday, RegTry.Gender)
	if err != nil { //did not exec query (syntax)
		return err
	}

	affected_rows, err := result.RowsAffected()
	if affected_rows <= 0 || err != nil { //Did not insert
		return err
	}
	return nil
}

//checkLogin checks if a user pass combination is valid for the specified auth try.
//
//Returns true if login is valid, false otherwise and report errors.
func checkLogin(AuthTry loginRequests.Auth) (bool, error) {
	var num_rows int
	var password_hash string
	var salt int

	stmtQuery, err := configurations.ArchivesPool.Prepare("SELECT COUNT(*) AS num_rows, HEX(hash_pwd), hash_salt FROM users WHERE username = ? GROUP BY hash_pwd, hash_salt")
	if err != nil {
		return false, err
	}
	defer stmtQuery.Close()

	result, err := stmtQuery.Query(AuthTry.Username)
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

//userInArchives checks into the archives if a user is already registered.
//
//Returns true if the user is already in the archives, false otherwise.
func isRegistered(username string, source_email string) (bool, error) {
	stmtQuery, err := configurations.ArchivesPool.Prepare("SELECT COUNT(*) AS num_rows FROM users WHERE username = ? OR email = ?")
	if err != nil { //cannot check, consider error
		return true, err
	}
	defer stmtQuery.Close()

	rows, err := stmtQuery.Query(username, source_email)
	if err != nil { //cannot check, consider error
		return true, err
	}
	defer stmtQuery.Close()

	var num_rows int
	rows.Scan(&num_rows)
	return num_rows > 0, nil
}

func updateCacheNewUserSession(email string) (string, error) {
	return controllerSharedFuncs.UpdateCacheNewSession(constants.LOGGED_USERS_SET, email, constants.CACHE_REFRESH_INTERVAL)
}
