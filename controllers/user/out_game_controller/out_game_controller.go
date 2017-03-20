package outGameController

import (
	"errors"
	"net/http"

	"sanino/gamemate/controllers/shared"
	"sanino/gamemate/controllers/user/session_controller"

	"sanino/gamemate/models/shared/responses/errors"
	"sanino/gamemate/models/user/requests/out_game"
	"sanino/gamemate/models/user/responses/out_game"

	"github.com/labstack/echo"
)

//HandleMyEnabledGamesForUser handles the request from a user to see his enabled games.
func HandleMyEnabledGamesForUser(context echo.Context) error {
	request := outGameRequests.UserGameList{}
	err := request.FromForm(context)
	if err != nil {
		errorResponse := errorResponses.ErrorDetail{}
		context.Logger().Print(err)
		errorResponse.FromError(err, http.StatusBadRequest)
		return context.JSON(http.StatusBadRequest, errorResponse)
	}
	if val, err := controllerSharedFuncs.IsValidAPI_Token(request.API_Token); !val || err != nil {
		errorResponse := errorResponses.ErrorDetail{}
		context.Logger().Print(errors.New("Rejected by the system, requestor not valid"))
		errorResponse.FromError(errors.New("Rejected by the system, request not valid"), http.StatusBadRequest)
		return context.JSON(http.StatusBadRequest, errorResponse)
	}
	userID, err := sessionController.GetUserIDFromSessionToken(request.SessionToken)
	if err != nil {
		errorResponse := errorResponses.ErrorDetail{}
		context.Logger().Print(errors.New("Rejected by the system, invalid session"))
		errorResponse.FromError(errors.New("Rejected by the system, invalid session"), http.StatusBadRequest)
		return context.JSON(http.StatusBadRequest, errorResponse)
	}
	games, cacheUpdated, err := GetEnabledGameIDs(userID)
	if err != nil {
		errorResponse := errorResponses.ErrorDetail{}
		context.Logger().Print(err)
		errorResponse.FromError(errors.New("Cannot get games"), http.StatusInternalServerError)
		return context.JSON(http.StatusInternalServerError, errorResponse)
	}
	if !cacheUpdated {
		context.Logger().Print("Request was successfull but cache was not updated")
	}

	responseFromServer := outGameResponses.MyEnabledGames{}
	responseFromServer.FromGameIDs(games)
	return context.JSON(http.StatusOK, &responseFromServer)
}

//HandleAllGamesForUser handles a request to show all games summarized data
//(e.g. name + ID + currentlyPlaying)
func HandleAllGamesForUser(context echo.Context) error {
	request := outGameRequests.UserGameList{}
	err := request.FromForm(context)
	if err != nil {
		errorResponse := errorResponses.ErrorDetail{}
		context.Logger().Print(err)
		errorResponse.FromError(err, http.StatusBadRequest)
		return context.JSON(http.StatusBadRequest, errorResponse)
	}
	if val, err := controllerSharedFuncs.IsValidAPI_Token(request.API_Token); !val || err != nil {
		errorResponse := errorResponses.ErrorDetail{}
		context.Logger().Print(errors.New("Rejected by the system, requestor not valid"))
		errorResponse.FromError(errors.New("Rejected by the system, request not valid"), http.StatusBadRequest)
		return context.JSON(http.StatusBadRequest, errorResponse)
	}
	userID, err := sessionController.GetUserIDFromSessionToken(request.SessionToken)
	if err != nil {
		errorResponse := errorResponses.ErrorDetail{}
		context.Logger().Print(errors.New("Rejected by the system, invalid session"))
		errorResponse.FromError(errors.New("Rejected by the system, invalid session"), http.StatusBadRequest)
		return context.JSON(http.StatusBadRequest, errorResponse)
	}
	games, cacheUpdated, err := GetGames(userID)
	if err != nil {
		errorResponse := errorResponses.ErrorDetail{}
		context.Logger().Print(err)
		errorResponse.FromError(errors.New("Cannot get games"), http.StatusInternalServerError)
		return context.JSON(http.StatusInternalServerError, errorResponse)
	}
	if !cacheUpdated {
		context.Logger().Print("Request was successfull but cache was not updated")
	}

	responseFromServer := outGameResponses.UserGameList{}
	responseFromServer.FromGames(games)
	return context.JSON(http.StatusOK, &responseFromServer)
}

func HandleGameDetail(context echo.Context) error {
	request := outGameRequests.UserGameDetail{}
	err := request.FromForm(context)
	if err != nil {
		errorResponse := errorResponses.ErrorDetail{}
		context.Logger().Print(err)
		errorResponse.FromError(err, http.StatusBadRequest)
		return context.JSON(http.StatusBadRequest, errorResponse)
	}
	val, err := controllerSharedFuncs.IsValidAPI_Token(request.API_Token)
	if !val || err != nil {
		errorResponse := errorResponses.ErrorDetail{}
		context.Logger().Print(errors.New("Rejected by the system, requestor not valid"))
		errorResponse.FromError(errors.New("Rejected by the system, request not valid"), http.StatusBadRequest)
		return context.JSON(http.StatusBadRequest, errorResponse)
	}
	userID, err := sessionController.GetUserIDFromSessionToken(request.SessionToken)
	if err != nil {
		errorResponse := errorResponses.ErrorDetail{}
		context.Logger().Print(errors.New("Rejected by the system, invalid session"))
		errorResponse.FromError(errors.New("Rejected by the system, invalid session"), http.StatusBadRequest)
		return context.JSON(http.StatusBadRequest, errorResponse)
	}
	game, cacheUpdated, err := getGameDetail(request.GameID, userID)
	if !cacheUpdated {
		context.Logger().Print("Game Detail : Cache not updated")
	}
	if err != nil {
		errorResponse := errorResponses.ErrorDetail{}
		context.Logger().Print("Error during game detail get : " + err.Error())
		errorResponse.FromError(errors.New("Cannot get game details"), http.StatusBadRequest)
		return context.JSON(http.StatusBadRequest, errorResponse)
	}
	response := outGameResponses.UserGameDetail{}
	response.FromGame(game)
	return context.JSON(http.StatusOK, &response)
}
