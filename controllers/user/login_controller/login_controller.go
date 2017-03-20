package loginController

//TODO : use controllerSharedFuncs.IsValidAPI_Token
import (
	"errors"
	"net/http"

	"sanino/gamemate/configurations"
	"sanino/gamemate/constants"
	"sanino/gamemate/controllers/shared"

	"sanino/gamemate/models/shared/responses/errors"
	"sanino/gamemate/models/user/requests/login"
	"sanino/gamemate/models/user/responses/login"

	"github.com/garyburd/redigo/redis"
	"github.com/labstack/echo"
)

//HandleAuth handles the authentication of the user for the system.
func HandleAuth(context echo.Context) error {
	errorResp := errorResponses.ErrorDetail{}
	var isLoggable bool
	AuthTry := loginRequests.Auth{}
	err := AuthTry.FromForm(context)

	if err != nil {
		context.Logger().Print(err)
		errorResp.FromError(errors.New("Bad Request"), http.StatusBadRequest)
		return context.JSON(http.StatusBadRequest, &errorResp)
	}

	isLoggable, userID, err := checkLogin(AuthTry)
	if err != nil {
		context.Logger().Print(err)
		errorResp.FromError(errors.New("Cannot Login User"), http.StatusInternalServerError)
		return context.JSON(http.StatusInternalServerError, &errorResp)
	}
	if !isLoggable {
		errMsg := "User - Pwd Combination not correct"
		errorResp.FromError(errors.New(errMsg), 1)
		context.Logger().Print(errMsg)
		return context.JSON(http.StatusBadRequest, &errorResp)
	}
	token, err := updateCacheNewUserSession(userID, AuthTry.Username)
	if err != nil {
		context.Logger().Print(err)
		errorResp.FromError(errors.New("Cannot Login User"), http.StatusInternalServerError)
		return context.JSON(http.StatusInternalServerError, &errorResp)
	}
	responseFromServer := loginResponses.Auth{}
	responseFromServer.FromToken(token)
	return context.JSON(http.StatusCreated, &responseFromServer)
}

//HandleRegistration handles the registration of a user for the system.
func HandleRegistration(context echo.Context) error {
	errorResp := errorResponses.ErrorDetail{}
	var RegTry = loginRequests.Registration{}
	var err = RegTry.FromForm(context)
	if err != nil {
		context.Logger().Print(err)
		errorResp.FromError(errors.New("Bad Request"), http.StatusBadRequest)
		return context.JSON(http.StatusBadRequest, &errorResp)
	}
	//check if user already exists
	isRegisteredUser, err := isRegistered(RegTry.Username)
	if err != nil {
		context.Logger().Printf("error username: %v", err)
		errorResp.FromError(errors.New("Cannot insert user"), http.StatusInternalServerError)
		return context.JSON(http.StatusInternalServerError, errorResp)
	}
	if isRegisteredUser {
		context.Logger().Print("The user is already registered into the system")
		errorResp.FromError(errors.New("User already registered"), 2)
		return context.JSON(http.StatusBadRequest, &errorResp)
	}
	//else query and if query successful add user also into cache and reply with session_token
	//generating random salt
	userID, err := insertIntoArchives(RegTry)
	if err != nil {
		context.Logger().Print(err)
		errorResp.FromError(errors.New("Cannot Insert User"), http.StatusInternalServerError)
		return context.JSON(http.StatusInternalServerError, errorResp)
	}

	token, err := updateCacheNewUserSession(userID, RegTry.Username)
	if err != nil {
		context.Logger().Print(err)
		errorResp.FromError(errors.New("Cannot Insert User"), http.StatusInternalServerError)
		return context.JSON(http.StatusInternalServerError, errorResp)
	}

	//finished, sending token to client
	responseFromServer := loginResponses.Auth{}
	responseFromServer.FromToken(token)
	return context.JSON(http.StatusCreated, &responseFromServer)
}

//GetUserIDFromSessionToken gets the user ID from session token.
//
//Returns error if not found in cache.
func GetUserIDFromSessionToken(token string) (int64, error) {
	return controllerSharedFuncs.GetIDFromSessionSet(constants.LOGGED_USERS_SET, token)
}

//GetConnectedUsers get the number of connected users.
func GetConnectedUsers() (int64, error) {
	conn := configurations.CachePool.Get()
	defer conn.Close()
	return redis.Int64(conn.Do("ZCARD", constants.LOGGED_USERS_SET))
}
