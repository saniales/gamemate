package outGameController

import (
	"errors"
	"net/http"

	"sanino/gamemate/controllers/developer"
	"sanino/gamemate/controllers/user/session_controller"

	"sanino/gamemate/models/shared/responses/errors"
	"sanino/gamemate/models/user/requests/out_game"
	"sanino/gamemate/models/user/responses/out_game"

	"github.com/labstack/echo"
)

//HandleMyEnabledGamesForUser handles the request from a user to see his enabled games.
func HandleMyEnabledGamesForUser(context echo.Context) error {
	request := outGameRequests.MyGames{}
	err := request.FromForm(context)
	if err != nil {
		errorResp := errorResponses.ErrorDetail{}
		context.Logger().Print(err)
		errorResp.FromError(err, http.StatusBadRequest)
		return context.JSON(http.StatusBadRequest, errorResp)
	}
	if val, err := developerController.IsValidAPI_Token(request.API_Token); !val || err != nil {
		errorResp := errorResponses.ErrorDetail{}
		context.Logger().Print(errors.New("Rejected by the system, requestor not valid"))
		errorResp.FromError(errors.New("Rejected by the system, request not valid"), http.StatusBadRequest)
		return context.JSON(http.StatusBadRequest, errorResp)
	}
	userID, err := sessionController.GetUserIDFromSessionToken(request.SessionToken)
	if err != nil {
		errorResp := errorResponses.ErrorDetail{}
		context.Logger().Print(errors.New("Rejected by the system, invalid session"))
		errorResp.FromError(errors.New("Rejected by the system, invalid session"), http.StatusBadRequest)
		return context.JSON(http.StatusBadRequest, errorResp)
	}
	games, cacheUpdated, err := GetEnabledGameIDs(userID)
	if err != nil {
		errorResp := errorResponses.ErrorDetail{}
		context.Logger().Print(err)
		errorResp.FromError(errors.New("Cannot get games"), http.StatusInternalServerError)
		return context.JSON(http.StatusInternalServerError, errorResp)
	}
	if !cacheUpdated {
		context.Logger().Print("Request was successfull but cache was not updated")
	}

	response := outGameResponses.MyEnabledGames{}
	response.FromGameIDs(games)
	return context.JSON(http.StatusOK, response)
}

//HandleAllGamesForUser handles a request to show all games summarized data
//(e.g. name + ID + currentlyPlaying)
func HandleAllGamesForUser(context echo.Context) error {
	request := outGameRequests.MyGames{}
	err := request.FromForm(context)
	if err != nil {
		errorResp := errorResponses.ErrorDetail{}
		context.Logger().Print(err)
		errorResp.FromError(err, http.StatusBadRequest)
		return context.JSON(http.StatusBadRequest, errorResp)
	}
	if val, err := developerController.IsValidAPI_Token(request.API_Token); !val || err != nil {
		errorResp := errorResponses.ErrorDetail{}
		context.Logger().Print(errors.New("Rejected by the system, requestor not valid"))
		errorResp.FromError(errors.New("Rejected by the system, request not valid"), http.StatusBadRequest)
		return context.JSON(http.StatusBadRequest, errorResp)
	}
	_, err = sessionController.GetUserIDFromSessionToken(request.SessionToken)
	if err != nil {
		errorResp := errorResponses.ErrorDetail{}
		context.Logger().Print(errors.New("Rejected by the system, invalid session"))
		errorResp.FromError(errors.New("Rejected by the system, invalid session"), http.StatusBadRequest)
		return context.JSON(http.StatusBadRequest, errorResp)
	}
	games, cacheUpdated, err := GetGames()
	if err != nil {
		errorResp := errorResponses.ErrorDetail{}
		context.Logger().Print(err)
		errorResp.FromError(errors.New("Cannot get games"), http.StatusInternalServerError)
		return context.JSON(http.StatusInternalServerError, errorResp)
	}
	if !cacheUpdated {
		context.Logger().Print("Request was successfull but cache was not updated")
	}

	response := outGameResponses.MyGames{}
	response.FromGames(games)
	return context.JSON(http.StatusOK, response)
}
