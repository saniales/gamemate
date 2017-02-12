package loginController

import (
	"errors"
	"net/http"
	"sanino/gamemate/constants"
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
		context.Logger().Print(err)
		errResp.FromError(errors.New("Bad Request"), http.StatusBadRequest)
		return context.JSON(http.StatusBadRequest, errResp)
	}

	isLoggable, err = checkLogin(AuthTry)
	if err != nil {
		context.Logger().Print(err)
		errResp.FromError(errors.New("Cannot Login User"), 500)
		return context.JSON(http.StatusInternalServerError, errResp)
	}
	if !isLoggable {
		errMsg := "Cannot login. User - Pwd Combination not correct"
		errResp.FromError(errors.New(errMsg), 1)
		context.Logger().Printf(errMsg)
		return context.JSON(http.StatusBadRequest, errResp)
	}
	halfHour, _ := time.ParseDuration("30m")
	token, err := updateCacheNewSession(AuthTry.Username, time.Duration(time.Now().Add(halfHour).UnixNano()))
	if err != nil {
		context.Logger().Print(err)
		errResp.FromError(errors.New("Cannot Login User"), 500)
		return context.JSON(http.StatusInternalServerError, errResp)
	}
	return context.JSON(http.StatusCreated, response.Auth{SessionToken: token})
}

//HandleRegistration handles the registration of a user for the system.
func HandleRegistration(context echo.Context) error {
	errResp := response.ErrorDetail{}
	var RegTry = request.Registration{}
	var err = RegTry.FromForm(context)
	if err != nil {
		context.Logger().Print(err)
		errResp.FromError(errors.New("Bad Request"), http.StatusBadRequest)
		return context.JSON(http.StatusBadRequest, errResp)
	}
	//check if user already exists
	isRegisteredUser, err := isRegistered(RegTry.Username, RegTry.Email)
	if err != nil {
		context.Logger().Print(err)
		errResp.FromError(errors.New("Cannot insert user"), 500)
		return context.JSON(http.StatusInternalServerError, errResp)
	}
	if isRegisteredUser {
		context.Logger().Print("The user " + RegTry.Username + " is already registered into the system")
		errResp.FromError(errors.New("User already registered"), 2)
		return context.JSON(http.StatusBadRequest, errResp)
	}
	//else query and if query successful add user also into cache and reply with session_token
	//generating random salt
	err = insertIntoArchives(RegTry)
	if err != nil {
		context.Logger().Print(err)
		errResp.FromError(errors.New("Cannot Insert User"), 500)
		return context.JSON(http.StatusInternalServerError, errResp)
	}

	token, err := updateCacheNewSession(RegTry.Username, constants.CACHE_REFRESH_INTERVAL)
	if err != nil {
		context.Logger().Print(err)
		errResp.FromError(errors.New("Cannot Insert User"), 500)
		return context.JSON(http.StatusInternalServerError, errResp)
	}

	//finished, sending token to client
	responseFromServer := response.Auth{SessionToken: token}
	return context.JSON(http.StatusCreated, responseFromServer)
}
