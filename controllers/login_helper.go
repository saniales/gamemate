package controllers

import (
	"crypto/sha512"
	"encoding/hex"
	"io"
	"log"
	"math/rand"
	"sanino/gamemate/configurations"
	"sanino/gamemate/constants"
	"sanino/gamemate/models/request"
	"strconv"
	"time"
)

func generateToken() string {
	SHA512hash := sha512.New()
	rand.Seed(time.Now().UnixNano())
	randomNum := strconv.Itoa(rand.Intn(constants.MAX_NUMBER_SALT))

	return string(SHA512hash.Sum([]byte(randomNum))[:])
}

func updateCacheNewSession(username string, expiration int64) (string, error) {
	token := generateToken()
	conn := configurations.CachePool.Get()

	err := conn.Send("MULTI")
	if err != nil {
		log.Fatal("cannot start redis transaction, quitting... Error detail =>" + err.Error())
		return constants.INVALID_TOKEN, err
	}

	err = conn.Send("ZADD", "logged_tokens", expiration, token) //sets the cache for the token expire 30 minutes
	if err != nil {
		log.Fatal("cannot send ZADD command, quitting... Error detail =>" + err.Error())
		return constants.INVALID_TOKEN, err
	}
	err = conn.Send("SET", "user/"+username+"/token", token, "EX", "1800")
	//if set when a user logons with an expired key it is removed from cache and set
	if err != nil {
		log.Fatal("cannot send SET user/.../token [v] command, quitting... Error detail =>" + err.Error())
		return constants.INVALID_TOKEN, err
	}

	err = conn.Send("SET", "token/"+token+"/user", username, "EX", "1800")
	//if set when a user logons with an expired key it is removed from cache and set
	if err != nil {
		log.Fatal("cannot send SET token/.../user [v] command, quitting... Error detail =>" + err.Error())
		return constants.INVALID_TOKEN, err
	}

	err = conn.Send("EXEC")
	if err != nil {
		log.Print("cannot commit transaction, quitting... Error detail =>" + err.Error())
		return constants.INVALID_TOKEN, err
	}
	return token, nil
}

//userInArchives checks into the archives if a user is already registered.
//
//Returns true if the user is already in the archives, false otherwise.
func isRegistered(username string) (bool, error) {
	stmtQuery, err := configurations.ArchivesPool.Prepare("SELECT COUNT(*) AS num_rows FROM user WHERE username = ?")
	if err != nil { //cannot check, consider error
		return true, err
	}
	defer stmtQuery.Close()

	rows, err := stmtQuery.Query(username)
	if err != nil { //cannot check, consider error
		log.Print(err)
		return true, err
	}

	var num_rows int
	rows.Scan(&num_rows)
	if num_rows > 0 { //user already registered
		log.Print(err)
		return true, nil
	}
	return false, nil
}

//insertIntoArchives (without check of previous insertions, only error reporting)
//inserts a new User into the archives, doing the salty & hashy work.
func insertIntoArchives(RegTry request.Registration) error {
	rand.Seed(time.Now().UTC().UnixNano())
	salt := rand.Intn(constants.MAX_NUMBER_SALT)
	saltedPass := RegTry.Password + strconv.Itoa(salt)
	stmtQuery, err := configurations.ArchivesPool.Prepare("INSERT INTO users (ID, username, password, salt, email, birthday, gender) VALUES (NULL, ?, SHA512(?), ?, ?, ?, ?)")
	if err != nil {
		log.Print(err)
		return err
	}
	defer stmtQuery.Close()

	result, err := stmtQuery.Exec(RegTry.Username, saltedPass, salt, RegTry.Email, RegTry.Birthday, RegTry.Gender)
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
func checkLogin(AuthTry request.Auth) (bool, error) {
	var num_rows int
	var password_hash string
	var salt int

	stmtQuery, err := configurations.ArchivesPool.Prepare("SELECT COUNT(*) AS num_rows, password, salt FROM users WHERE username = ?")
	if err != nil {
		log.Print(err)
		return false, err
	}

	result, err := stmtQuery.Query(AuthTry.Username)
	if err != nil {
		log.Print(err)
		return false, err
	}

	err = result.Scan(&num_rows, &password_hash, &salt)
	if err != nil {
		log.Print(err)
		return false, err
	}
	hash := sha512.New()
	salted_pwd := AuthTry.Password + strconv.Itoa(salt)
	io.WriteString(hash, salted_pwd)
	salted_hash := hex.EncodeToString(hash.Sum(nil))
	return salted_hash == password_hash, nil
}
