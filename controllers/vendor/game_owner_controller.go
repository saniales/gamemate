package gameOwnerController

import (
	"errors"
	"fmt"
	"net/http"

	"sanino/gamemate/controllers/developer"
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
		return context.JSON(http.StatusBadRequest, errorResp)
	}
	if val, err := developerController.IsValidAPI_Token(request.API_Token); !val || err != nil {
		errorResp := errorResponses.ErrorDetail{}
		context.Logger().Print(errors.New("Rejected by the system, requestor not valid"))
		errorResp.FromError(err, http.StatusBadRequest)
		return context.JSON(http.StatusBadRequest, errorResp)
	}
	email, err := getOwnerEmailFromSessionToken(request.SessionToken)
	if err != nil {
		errorResp := errorResponses.ErrorDetail{}
		context.Logger().Print(fmt.Errorf("%s token rejected by the system, Invalid Session", request.SessionToken))
		errorResp.FromError(errors.New("Rejected by the system"), http.StatusBadRequest)
		return context.JSON(http.StatusBadRequest, errorResp)
	}
	token, err := addGameInArchives(email)
	if err != nil {
		errorResp := errorResponses.ErrorDetail{}
		context.Logger().Print(fmt.Errorf("Cannot create new API Token, error => %v", err))
		errorResp.FromError(errors.New("Cannot create API Token"), http.StatusInternalServerError)
		return context.JSON(http.StatusInternalServerError, errorResp)
	}
	response := gameOwnerResponses.AddGame{}

	panic("IMPLEMENT IT, FAG")
	return context.JSON(http.StatusCreated, response)
}

//HandleDropAPI_Token handles a request to remove a developer API Token.
func HandleDropAPI_Token(context echo.Context) error {
	request := developerRequests.DropToken{}
	err := request.FromForm(context)
	if err != nil {
		errorResp := errorResponses.ErrorDetail{}
		errorResp.FromError(err, http.StatusBadRequest)
		return context.JSON(http.StatusBadRequest, errorResp)
	}

	IsValid, err := IsValidAPI_Token(request.API_Token)
	if !IsValid || err != nil {
		context.Logger().Print(fmt.Errorf("API Token %s rejected", request.API_Token))
		errorResp := errorResponses.ErrorDetail{}
		errorResp.FromError(errors.New("Rejected by the system"), http.StatusBadRequest)
		return context.JSON(http.StatusBadRequest, errorResp)
	}

	err = removeAPI_TokenFromCache(request.TokenToDrop)
	if err != nil {
		context.Logger().Print(fmt.Errorf("%s API Token not removed. Error => %v", request.TokenToDrop, err))
		errorResp := errorResponses.ErrorDetail{}
		errorResp.FromError(errors.New("Cannot remove API Token"), http.StatusInternalServerError)
		return context.JSON(http.StatusInternalServerError, errorResp)
	}

	err = removeAPI_TokenFromArchives(request.TokenToDrop)
	if err != nil {
		context.Logger().Print(fmt.Errorf("%s API Token not removed. Error => %v", request.TokenToDrop, err))
		errorResp := errorResponses.ErrorDetail{}
		errorResp.FromError(errors.New("Cannot remove API Token"), http.StatusInternalServerError)
		return context.JSON(http.StatusInternalServerError, errorResp)
	}

	response := developerResponses.DropToken{}
	response.FromOldAPIToken(request.TokenToDrop)
	return context.JSON(http.StatusOK, response)
}

//HandleRegistration handles a request to register a developer.
func HandleRegistration(context echo.Context) error {
	request := developerRequests.DevRegistration{}
	err := request.FromForm(context)
	if err != nil {
		errorResp := errorResponses.ErrorDetail{}
		errorResp.FromError(err, http.StatusBadRequest)
		return context.JSON(http.StatusBadRequest, errorResp)
	}

	err = registerDeveloper(request)
	if err != nil {
		errorResp := errorResponses.ErrorDetail{}
		errorResp.FromError(err, http.StatusInternalServerError)
		return context.JSON(http.StatusBadRequest, errorResp)
	}

	token, err := updateCacheWithSessionDeveloperToken(request.Email)
	if err != nil {
		errorResp := errorResponses.ErrorDetail{}
		errorResp.FromError(errors.New("User registered, but I did not login automatically, try to login later"), http.StatusBadRequest)
		return context.JSON(http.StatusInternalServerError, errorResp)
	}

	responseFromServer := developerResponses.DevAuth{}
	responseFromServer.FromToken(token)
	return context.JSON(http.StatusCreated, responseFromServer)
}

//HandleLogin handles login requests for developers.
func HandleLogin(context echo.Context) error {
	request := developerRequests.DevAuth{}
	err := request.FromForm(context)
	if err != nil {
		errorResp := errorResponses.ErrorDetail{}
		errorResp.FromError(err, http.StatusBadRequest)
		return context.JSON(http.StatusBadRequest, errorResp)
	}

	IsValid, err := IsValidAPI_Token(request.API_Token)
	if !IsValid || err != nil {
		context.Logger().Print(fmt.Errorf("API Token %s rejected", request.API_Token))
		errorResp := errorResponses.ErrorDetail{}
		errorResp.FromError(errors.New("Rejected by the system"), http.StatusBadRequest)
		return context.JSON(http.StatusBadRequest, errorResp)
	}

	isLoggable, err := checkLogin(request)
	if err != nil {
		errorResp := errorResponses.ErrorDetail{}
		context.Logger().Print(err)
		errorResp.FromError(errors.New("Login failed"), http.StatusBadRequest)
		return context.JSON(http.StatusBadRequest, errorResp)
	}
	if !isLoggable {
		errorResp := errorResponses.ErrorDetail{}
		context.Logger().Print(err)
		errorResp.FromError(errors.New("User - Password combination wrong, retry"), http.StatusBadRequest)
		return context.JSON(http.StatusBadRequest, errorResp)
	}
	token, err := updateCacheWithSessionDeveloperToken(request.Email)
	if err != nil {
		errorResp := errorResponses.ErrorDetail{}
		context.Logger().Print(err)
		errorResp.FromError(errors.New("Temporary error, retry in a few seconds"), http.StatusInternalServerError)
		return context.JSON(http.StatusInternalServerError, errorResp)
	}
	response := developerResponses.DevAuth{}
	response.FromToken(token)
	return context.JSON(http.StatusCreated, response)
}
