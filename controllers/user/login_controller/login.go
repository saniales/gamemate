package loginController

import (
	"errors"
	"net/http"

	"sanino/gamemate/models/shared/responses/errors"
	"sanino/gamemate/models/user/requests/login"
	"sanino/gamemate/models/user/responses/login"

	"github.com/labstack/echo"
)

//HandleAuth handles the authentication of the user for the system.
func HandleAuth(context echo.Context) error {
	errResp := errorResponses.ErrorDetail{}
	var isLoggable bool
	var AuthTry = loginRequests.Auth{}
	var err = AuthTry.FromForm(context)

	if err != nil {
		context.Logger().Print(err)
		errResp.FromError(errors.New("Bad Request"), http.StatusBadRequest)
		return context.JSON(http.StatusBadRequest, errResp)
	}

	isLoggable, userID, err := checkLogin(AuthTry) //TODO: doubt, should i return an "User" struct??
	if err != nil {
		context.Logger().Print(err)
		errResp.FromError(errors.New("Cannot Login User"), http.StatusInternalServerError)
		return context.JSON(http.StatusInternalServerError, errResp)
	}
	if !isLoggable {
		errMsg := "Cannot login. User - Pwd Combination not correct"
		errResp.FromError(errors.New(errMsg), 1)
		context.Logger().Printf(errMsg)
		return context.JSON(http.StatusBadRequest, errResp)
	}
	token, err := updateCacheNewUserSession(userID, AuthTry.Username)
	if err != nil {
		context.Logger().Print(err)
		errResp.FromError(errors.New("Cannot Login User"), http.StatusInternalServerError)
		return context.JSON(http.StatusInternalServerError, errResp)
	}
	responseFromServer := loginResponses.Auth{}
	responseFromServer.FromToken(token)
	return context.JSON(http.StatusCreated, responseFromServer)
}

//HandleRegistration handles the registration of a user for the system.
func HandleRegistration(context echo.Context) error {
	errResp := errorResponses.ErrorDetail{}
	var RegTry = loginRequests.Registration{}
	var err = RegTry.FromForm(context)
	if err != nil {
		context.Logger().Print(err)
		errResp.FromError(errors.New("Bad Request"), http.StatusBadRequest)
		return context.JSON(http.StatusBadRequest, errResp)
	}
	//check if user already exists
	isRegisteredUser, err := isRegistered(RegTry.Username)
	//NOTE: doubl connection to DB, not so efficient, replace with a boolean
	//combination to avoid second call.
	isRegisteredUserEmail, errMail := isRegistered(RegTry.Email)
	if err != nil || errMail != nil {
		context.Logger().Printf("error username: %v, error mail:%v", err, errMail)
		errResp.FromError(errors.New("Cannot insert user"), http.StatusInternalServerError)
		return context.JSON(http.StatusInternalServerError, errResp)
	}
	if isRegisteredUser || isRegisteredUserEmail {
		context.Logger().Print("The user is already registered into the system")
		errResp.FromError(errors.New("User already registered"), 2)
		return context.JSON(http.StatusBadRequest, errResp)
	}
	//else query and if query successful add user also into cache and reply with session_token
	//generating random salt
	userID, err := insertIntoArchives(RegTry)
	if err != nil {
		context.Logger().Print(err)
		errResp.FromError(errors.New("Cannot Insert User"), http.StatusInternalServerError)
		return context.JSON(http.StatusInternalServerError, errResp)
	}

	token, err := updateCacheNewUserSession(userID, RegTry.Username)
	if err != nil {
		context.Logger().Print(err)
		errResp.FromError(errors.New("Cannot Insert User"), http.StatusInternalServerError)
		return context.JSON(http.StatusInternalServerError, errResp)
	}

	//finished, sending token to client
	responseFromServer := loginResponses.Auth{}
	responseFromServer.FromToken(token)
	return context.JSON(http.StatusCreated, responseFromServer)
}
