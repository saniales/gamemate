package controllers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"sanino/gamemate/constants"
	"sanino/gamemate/models/request"
	"sanino/gamemate/models/response"
	"time"

	"github.com/labstack/echo"
)

//HandleAuth handles the authentication of the user for the system.
func HandleAuth(c echo.Context) error {
	var isLoggable bool
	var AuthTry = request.Auth{}
	var err = AuthTry.FromForm(c)

	if err != nil {
		return err
	}

	isLoggable, err = checkLogin(AuthTry)
	if err != nil {
		log.Print(err)
		return err
	}
	if !isLoggable {
		errMsg := "Cannot log [ %v ] auth try. User - Pwd Combination not correct"
		log.Printf(errMsg, AuthTry)
		return fmt.Errorf(errMsg, AuthTry)
	}
	// answers with session token valid for 30 minutes from last session
	// it has to be created and put in redis after correct authentication.
	// otherwise the system must reply with a system_error struct {Code : 1}
	// if debug it includes a message {errCode : 1, message : "ZIOBANANA"}
	return nil
}

//HandleRegistration handles the registration of a user for the system.
func HandleRegistration(context echo.Context) error {
	var RegTry = request.Registration{}
	var err = RegTry.FromForm(context)
	if err != nil {
		log.Print(err)
		return err
	}
	//check if user already exists
	isRegisteredUser, err := isRegistered(RegTry.Username)
	if err != nil {
		log.Print(err)
		return err
	}
	if isRegisteredUser {
		log.Print("The user " + RegTry.Username + " is already registered into the system, reporting error...")
		return errors.New("This user is already registered")
	}
	//else query and if query successful add user also into cache and reply with session_token
	//generating random salt
	err = insertIntoArchives(RegTry)
	if err != nil {
		log.Print(err)
		return err
	}

	span, _ := time.ParseDuration(constants.MAX_DURATION)
	token, err := updateCacheNewSession(time.Now().Add(span).UnixNano())
	if err != nil {
		log.Print(err)
	}
	//finished, sending token to client
	responseFromServer := response.Auth{SessionToken: token}
	return context.JSON(http.StatusCreated, responseFromServer)
}
