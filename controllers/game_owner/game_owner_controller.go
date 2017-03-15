package gameOwnerController

import (
	"errors"
	"fmt"
	"net/http"

	"sanino/gamemate/controllers/shared"
	"sanino/gamemate/controllers/user/login_controller"

	"sanino/gamemate/models/game_owner/requests"
	"sanino/gamemate/models/game_owner/responses"
	"sanino/gamemate/models/shared/responses/errors"

	"github.com/labstack/echo"
)

//HandleAddGame handles a request to add a developer API Token.
func HandleAddGame(context echo.Context) error {
	request := gameOwnerRequests.AddGame{}
	err := request.FromForm(context)
	if err != nil {
		errorResp := errorResponses.ErrorDetail{}
		context.Logger().Print(err)
		errorResp.FromError(err, http.StatusBadRequest)
		return context.JSON(http.StatusBadRequest, &errorResp)
	}
	if val, err := controllerSharedFuncs.IsValidAPI_Token(request.API_Token); !val || err != nil {
		errorResp := errorResponses.ErrorDetail{}
		context.Logger().Print(errors.New("Rejected by the system, requestor not valid"))
		errorResp.FromError(err, http.StatusBadRequest)
		return context.JSON(http.StatusBadRequest, &errorResp)
	}
	ownerID, err := getOwnerIDFromSessionToken(request.SessionToken)
	if err != nil {
		errorResp := errorResponses.ErrorDetail{}
		context.Logger().Print(fmt.Errorf("%s token rejected by the system, Invalid Session", request.SessionToken))
		errorResp.FromError(errors.New("Rejected by the system"), http.StatusBadRequest)
		return context.JSON(http.StatusBadRequest, &errorResp)
	}
	gameID, err := addGameInArchives(ownerID, request.GameName, request.GameDescription, request.MaxPlayers)
	if err != nil {
		errorResp := errorResponses.ErrorDetail{}
		context.Logger().Print(fmt.Errorf("Cannot create new API Token, error => %v", err))
		errorResp.FromError(errors.New("Cannot create API Token"), http.StatusInternalServerError)
		return context.JSON(http.StatusInternalServerError, &errorResp)
	}
	response := gameOwnerResponses.AddGame{}

	response.FromGameID(gameID)
	return context.JSON(http.StatusCreated, &response)
}

//HandleRemoveGame handles a request to remove a developer API Token.
func HandleRemoveGame(context echo.Context) error {
	request := gameOwnerRequests.RemoveGame{}
	err := request.FromForm(context)
	if err != nil {
		errorResp := errorResponses.ErrorDetail{}
		errorResp.FromError(err, http.StatusBadRequest)
		return context.JSON(http.StatusBadRequest, &errorResp)
	}

	IsValid, err := controllerSharedFuncs.IsValidAPI_Token(request.API_Token)
	if !IsValid || err != nil {
		context.Logger().Print(fmt.Errorf("API Token %s rejected", request.API_Token))
		errorResp := errorResponses.ErrorDetail{}
		errorResp.FromError(errors.New("Rejected by the system"), http.StatusBadRequest)
		return context.JSON(http.StatusBadRequest, &errorResp)
	}

	ownerID, err := getOwnerIDFromSessionToken(request.SessionToken)
	if err != nil {
		errorResp := errorResponses.ErrorDetail{}
		context.Logger().Print(fmt.Errorf("%s token rejected by the system, Invalid Session", request.SessionToken))
		errorResp.FromError(errors.New("Rejected by the system"), http.StatusBadRequest)
		return context.JSON(http.StatusBadRequest, &errorResp)
	}

	err = removeGameFromCache(request.GameID)
	if err != nil {
		context.Logger().Print(fmt.Errorf("Game with ID:%d not removed. Error => %v", request.GameID, err))
		errorResp := errorResponses.ErrorDetail{}
		errorResp.FromError(errors.New("Cannot remove API Token"), http.StatusInternalServerError)
		return context.JSON(http.StatusInternalServerError, &errorResp)
	}

	err = removeGameFromArchives(ownerID, request.GameID)
	if err != nil {
		context.Logger().Print(fmt.Errorf("Game with ID:%d not removed. Error => %v", request.GameID, err))
		errorResp := errorResponses.ErrorDetail{}
		errorResp.FromError(errors.New("Cannot remove API Token"), http.StatusInternalServerError)
		return context.JSON(http.StatusInternalServerError, &errorResp)
	}

	response := gameOwnerResponses.RemoveGame{}
	response.FromGameID(request.GameID)
	return context.JSON(http.StatusOK, &response)
}

//HandleRegistration handles a request to register a developer.
func HandleRegistration(context echo.Context) error {
	request := gameOwnerRequests.GameOwnerRegistration{}
	err := request.FromForm(context)
	if err != nil {
		errorResp := errorResponses.ErrorDetail{}
		errorResp.FromError(err, http.StatusBadRequest)
		return context.JSON(http.StatusBadRequest, &errorResp)
	}

	IsValid, err := controllerSharedFuncs.IsValidAPI_Token(request.API_Token)
	if !IsValid || err != nil {
		context.Logger().Print(fmt.Errorf("API Token %s rejected", request.API_Token))
		errorResp := errorResponses.ErrorDetail{}
		errorResp.FromError(errors.New("Rejected by the system"), http.StatusBadRequest)
		return context.JSON(http.StatusBadRequest, &errorResp)
	}

	ownerID, err := registerOwner(request)
	if err != nil {
		errorResp := errorResponses.ErrorDetail{}
		errorResp.FromError(err, http.StatusInternalServerError)
		return context.JSON(http.StatusBadRequest, &errorResp)
	}

	token, err := updateCacheWithSessionOwnerToken(ownerID)
	if err != nil {
		errorResp := errorResponses.ErrorDetail{}
		errorResp.FromError(errors.New("User registered, but I did not login automatically, try to login later"), http.StatusBadRequest)
		return context.JSON(http.StatusInternalServerError, &errorResp)
	}

	responseFromServer := gameOwnerResponses.GameOwnerAuth{}
	responseFromServer.FromToken(token)
	return context.JSON(http.StatusCreated, responseFromServer)
}

//HandleLogin handles login requests for developers.
func HandleLogin(context echo.Context) error {
	request := gameOwnerRequests.GameOwnerAuth{}
	err := request.FromForm(context)
	if err != nil {
		errorResp := errorResponses.ErrorDetail{}
		errorResp.FromError(err, http.StatusBadRequest)
		return context.JSON(http.StatusBadRequest, &errorResp)
	}

	IsValid, err := controllerSharedFuncs.IsValidAPI_Token(request.API_Token)
	if !IsValid || err != nil {
		context.Logger().Print(fmt.Errorf("API Token %s rejected", request.API_Token))
		errorResp := errorResponses.ErrorDetail{}
		errorResp.FromError(errors.New("Rejected by the system"), http.StatusBadRequest)
		return context.JSON(http.StatusBadRequest, &errorResp)
	}

	isLoggable, ownerID, err := checkLogin(request)
	if err != nil {
		errorResp := errorResponses.ErrorDetail{}
		context.Logger().Print(err)
		errorResp.FromError(errors.New("Login failed"), http.StatusBadRequest)
		return context.JSON(http.StatusBadRequest, &errorResp)
	}
	if !isLoggable {
		errorResp := errorResponses.ErrorDetail{}
		context.Logger().Print(err)
		errorResp.FromError(errors.New("User - Password combination wrong, retry"), http.StatusBadRequest)
		return context.JSON(http.StatusBadRequest, &errorResp)
	}
	token, err := updateCacheWithSessionOwnerToken(ownerID)
	if err != nil {
		errorResp := errorResponses.ErrorDetail{}
		context.Logger().Print(err)
		errorResp.FromError(errors.New("Temporary error, retry in a few seconds"), http.StatusInternalServerError)
		return context.JSON(http.StatusInternalServerError, &errorResp)
	}
	response := gameOwnerResponses.GameOwnerAuth{}
	response.FromToken(token)
	return context.JSON(http.StatusCreated, &response)
}

//HandleGameAction handles the requests to enable/disable a game.
func HandleGameAction(context echo.Context) error {
	request := gameOwnerRequests.GameOwnerAction{}
	err := request.FromForm(context)
	if err != nil {
		errorResp := errorResponses.ErrorDetail{}
		errorResp.FromError(err, http.StatusBadRequest)
		return context.JSON(http.StatusBadRequest, &errorResp)
	}

	IsValid, err := controllerSharedFuncs.IsValidAPI_Token(request.API_Token)
	if !IsValid || err != nil {
		context.Logger().Print(fmt.Errorf("API Token %s rejected", request.API_Token))
		errorResp := errorResponses.ErrorDetail{}
		errorResp.FromError(errors.New("Rejected by the system"), http.StatusBadRequest)
		return context.JSON(http.StatusBadRequest, &errorResp)
	}

	//requires strictly seller/user
	ownerID, errOwner := getOwnerIDFromSessionToken(request.SessionToken)
	if errOwner != nil {
		userID, errUser := loginController.GetUserIDFromSessionToken(request.SessionToken)
		if errUser != nil {
			errorResp := errorResponses.ErrorDetail{}
			context.Logger().Print(fmt.Errorf("%s token rejected by the system, Invalid Session", request.SessionToken))
			errorResp.FromError(errors.New("Rejected by the system"), http.StatusBadRequest)
			return context.JSON(http.StatusBadRequest, &errorResp)
		} else {
			if userID == request.UserID { //OK
				cacheUpdated, err := EnableDisableGameForUser(request.UserID, request.GameID, request.Action)
				if err != nil {
					errorResp := errorResponses.ErrorDetail{}
					context.Logger().Print(fmt.Errorf("Error in archives : %v", err))
					errorResp.FromError(errors.New("Cannot satisfy request"), http.StatusInternalServerError)
					return context.JSON(http.StatusInternalServerError, &errorResp)
				}
				if !cacheUpdated {
					context.Logger().Print("Game Action completed on archives, but not on cache")
				}
				//NOTE:OK!!!_________________________________________________
				response := gameOwnerResponses.GameOwnerAction{}
				response.FromGameID(request.GameID)
				return context.JSON(http.StatusOK, &response)
			}
			errorResp := errorResponses.ErrorDetail{}
			context.Logger().Print(fmt.Errorf("Request rejected by the system, Invalid Request : requestor ID invalid => %d intead of %d", request.UserID, userID))
			errorResp.FromError(errors.New("You don't have the permission to perform this action"), http.StatusBadRequest)
			return context.JSON(http.StatusBadRequest, &errorResp)
		}
	} else {
		//verify owner act on his games.
		ownerOfGame, cacheUpdated, err := GetOwnerOfGame(request.GameID)
		if err != nil {
			context.Logger().Print(fmt.Errorf("Enable/disable %v: Cannot satisfy request, query error", request))
			errorResp := errorResponses.ErrorDetail{}
			errorResp.FromError(errors.New("Cannot satisfy request"), http.StatusInternalServerError)
			return context.JSON(http.StatusInternalServerError, &errorResp)
		}
		if !cacheUpdated {
			context.Logger().Print("Request satisfied, but cache has not been updated")
		}
		if ownerOfGame != ownerID {
			context.Logger().Print(fmt.Errorf("Enable/disable %v: Cannot satisfy request, rejected owner", request))
			errorResp := errorResponses.ErrorDetail{}
			errorResp.FromError(errors.New("You don't have the permission to perform this action"), http.StatusBadRequest)
			return context.JSON(http.StatusBadRequest, &errorResp)
		}
	}
	cacheUpdated, err := EnableDisableGameForUser(request.UserID, request.GameID, request.Action)
	if err != nil {
		context.Logger().Print(fmt.Errorf("Enable/disable %v: Cannot satisfy request, query error", request))
		errorResp := errorResponses.ErrorDetail{}
		errorResp.FromError(errors.New("Cannot satisfy request"), http.StatusInternalServerError)
		return context.JSON(http.StatusInternalServerError, &errorResp)
	}
	if !cacheUpdated {
		context.Logger().Print("Game Action completed on archives, but not on cache")
	}
	response := gameOwnerResponses.GameOwnerAction{}
	response.FromGameID(request.GameID)
	return context.JSON(http.StatusOK, &response)
}

//HandleShowMyGames handles the request to show the games owned by a game_owner.
func HandleShowMyGames(context echo.Context) error {
	request := gameOwnerRequests.GameOwnerGameList{}
	err := request.FromForm(context)
	if err != nil {
		errorResp := errorResponses.ErrorDetail{}
		errorResp.FromError(err, http.StatusBadRequest)
		return context.JSON(http.StatusBadRequest, &errorResp)
	}

	IsValid, err := controllerSharedFuncs.IsValidAPI_Token(request.API_Token)
	if !IsValid || err != nil {
		context.Logger().Print(fmt.Errorf("API Token %s rejected", request.API_Token))
		errorResp := errorResponses.ErrorDetail{}
		errorResp.FromError(errors.New("Rejected by the system"), http.StatusBadRequest)
		return context.JSON(http.StatusBadRequest, &errorResp)
	}

	//requires strictly seller/user
	ownerID, errOwner := getOwnerIDFromSessionToken(request.SessionToken)
	if errOwner != nil {
		errorResp := errorResponses.ErrorDetail{}
		context.Logger().Print(fmt.Errorf("%s token rejected by the system, Invalid Session", request.SessionToken))
		errorResp.FromError(errors.New("Rejected by the system"), http.StatusBadRequest)
		return context.JSON(http.StatusBadRequest, &errorResp)
	}
	//games, cacheOk, err := GetGames(ownerID) TODO
	games, err := GetGames(ownerID)
	if err != nil {
		context.Logger().Printf("GameList error : Cannot satisfy request, error => %v", err)
		errorResp := errorResponses.ErrorDetail{}
		errorResp.FromError(errors.New("Cannot satisfy request"), http.StatusInternalServerError)
		return context.JSON(http.StatusInternalServerError, &errorResp)
	}
	response := gameOwnerResponses.GameOwnerGameList{}
	response.FromGames(games)
	return context.JSON(http.StatusOK, &response)
}
