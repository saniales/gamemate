package controllers

import (
	"errors"
	"log"
	"net/http"
	"sanino/gamemate/models/request"
	"sanino/gamemate/models/response"
	"time"

	"github.com/labstack/echo"
)

//HandleAuth handles the authentication of the user for the system.
func HandleAuth(context echo.Context) error {
	errResp := response.ErrorDetail{}
	var isLoggable bool
	var AuthTry = request.Auth{}
	var err = AuthTry.FromForm(context)

	if err != nil {
		log.Print(err)
		errResp.FromError(errors.New("Bad Request"), http.StatusBadRequest)
		return context.JSON(http.StatusBadRequest, errResp)
	}

	isLoggable, err = checkLogin(AuthTry)
	if err != nil {
		log.Print(err)
		errResp.FromError(errors.New("Internal server error"), 500)
		return context.JSON(http.StatusInternalServerError, errResp)
	}
	if !isLoggable {
		errMsg := "Cannot login. User - Pwd Combination not correct"
		errResp.FromError(errors.New(errMsg), 1)
		log.Printf(errMsg, AuthTry)
		return context.JSON(http.StatusBadRequest, errResp)
	}
	halfHour, _ := time.ParseDuration("30m")
	token, err := updateCacheNewSession(AuthTry.Username, time.Now().Add(halfHour).UnixNano())
	if err != nil {
		log.Print(err)
		errResp.FromError(errors.New("Internal Server Error"), 500)
		return context.JSON(http.StatusInternalServerError, errResp)
	}
	// answers with session token valid for 30 minutes from last session
	// it has to be created and put in redis after correct authentication.
	// otherwise the system must reply with a system_error struct {Code : 1}
	// if debug it includes a message {errCode : 1, message : "ZIOBANANA"}
	return context.JSON(http.StatusCreated, response.Auth{SessionToken: token})
}

//HandleRegistration handles the registration of a user for the system.
func HandleRegistration(context echo.Context) error {
	errResp := response.ErrorDetail{}
	var RegTry = request.Registration{}
	var err = RegTry.FromForm(context)
	if err != nil {
		log.Print(err)
		errResp.FromError(errors.New("Internal Server Error"), 500)
		return context.JSON(http.StatusInternalServerError, errResp)
	}
	//check if user already exists
	isRegisteredUser, err := isRegistered(RegTry.Username)
	if err != nil {
		log.Print(err)
		errResp.FromError(errors.New("Internal Server Error"), 500)
		return context.JSON(http.StatusInternalServerError, errResp)
	}
	if isRegisteredUser {
		log.Print("The user " + RegTry.Username + " is already registered into the system, reporting error...")
		errResp.FromError(errors.New("User already registered"), 2)
		return context.JSON(http.StatusInternalServerError, errResp)
	}
	//else query and if query successful add user also into cache and reply with session_token
	//generating random salt
	err = insertIntoArchives(RegTry)
	if err != nil {
		log.Print(err)
		errResp.FromError(errors.New("Internal Server Error"), 500)
		return context.JSON(http.StatusInternalServerError, errResp)
	}

	halfHour, _ := time.ParseDuration("30m")
	token, err := updateCacheNewSession(RegTry.Username, time.Now().Add(halfHour).UnixNano())
	if err != nil {
		log.Print(err)
		errResp.FromError(errors.New("Internal Server Error"), 500)
		return context.JSON(http.StatusInternalServerError, errResp)
	}

	//finished, sending token to client
	responseFromServer := response.Auth{SessionToken: token}
	return context.JSON(http.StatusCreated, responseFromServer)
}
