package developerController

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"sanino/gamemate/controllers/shared"

	"sanino/gamemate/models/developer/requests"
	"sanino/gamemate/models/developer/responses"
	"sanino/gamemate/models/shared/responses/errors"

	"github.com/labstack/echo"
)

func HandleAllTokensForDeveloper(context echo.Context) error {
	request := new(developerRequests.TokenList)
	err := request.FromForm(context)
	if err != nil {
		errorResp := errorResponses.ErrorDetail{}
		context.Logger().Print(err)
		errorResp.FromError(err, http.StatusBadRequest)
		errorResp.ErrorMessage += fmt.Sprintf("%v", context.Request())
		fmt.Print(errorResp.ErrorMessage)
		return context.JSON(http.StatusBadRequest, &errorResp)
	}
	if val, err := controllerSharedFuncs.IsValidAPI_Token(request.API_Token); !val || err != nil {
		errorResp := errorResponses.ErrorDetail{}
		context.Logger().Print(errors.New("Rejected by the system, requestor not valid"))
		errorResp.FromError(err, http.StatusBadRequest)
		return context.JSON(http.StatusBadRequest, &errorResp)
	}
	ID, err := getDevIDFromSessionToken(request.SessionToken)
	if err != nil {
		errorResp := errorResponses.ErrorDetail{}
		context.Logger().Print(fmt.Errorf("%s token rejected by the system, Invalid Session : error => %s", request.SessionToken, err.Error()))
		errorResp.FromError(errors.New("Rejected by the system"), http.StatusBadRequest)
		return context.JSON(http.StatusBadRequest, &errorResp)
	}
	Tokens, cacheUpdated, err := getAPITokensOfDeveloper(ID)
	if cacheUpdated {
		context.Logger().Print("Cache Updated with token_list of developer " + strconv.FormatInt(ID, 10))
	}
	if err != nil {
		context.Logger().Print(err.Error() + ", developerID : " + strconv.FormatInt(ID, 10))
		errorResp := errorResponses.ErrorDetail{}
		errorResp.FromError(errors.New("Cannot Get Tokens"), http.StatusInternalServerError)
		return context.JSON(http.StatusInternalServerError, &errorResp)
	}
	response := developerResponses.TokenList{}
	response.FromTokens(Tokens)
	return context.JSON(http.StatusOK, &response)
}

//HandleAddAPI_Token handles a request to add a developer API Token.
func HandleAddAPI_Token(context echo.Context) error {
	request := developerRequests.AddToken{}
	err := request.FromForm(context)
	if err != nil {
		errorResp := errorResponses.ErrorDetail{}
		context.Logger().Print(err)
		errorResp.FromError(err, http.StatusBadRequest)
		errorResp.ErrorMessage += fmt.Sprintf("%v", context.Request())
		fmt.Print(errorResp.ErrorMessage)
		return context.JSON(http.StatusBadRequest, &errorResp)
	}
	val, err := controllerSharedFuncs.IsValidAPI_Token(request.API_Token)
	if !val || err != nil {
		errorResp := errorResponses.ErrorDetail{}
		context.Logger().Print(errors.New("Rejected by the system, requestor not valid"))
		errorResp.FromError(err, http.StatusBadRequest)
		return context.JSON(http.StatusBadRequest, &errorResp)
	}
	ID, err := getDevIDFromSessionToken(request.SessionToken)
	if err != nil {
		errorResp := errorResponses.ErrorDetail{}
		context.Logger().Print(fmt.Errorf("%s token rejected by the system, Invalid Session", request.SessionToken))
		errorResp.FromError(errors.New("Rejected by the system"), http.StatusBadRequest)
		return context.JSON(http.StatusBadRequest, &errorResp)
	}
	token, err := addAPI_TokenInArchives(ID)
	if err != nil {
		errorResp := errorResponses.ErrorDetail{}
		context.Logger().Print(fmt.Errorf("Cannot create new API Token, error => %v", err))
		errorResp.FromError(errors.New("Cannot create API Token"), http.StatusInternalServerError)
		return context.JSON(http.StatusInternalServerError, errorResp)
	}
	err = controllerSharedFuncs.UpdateCacheWithAPI_Token(token)
	if err != nil {
		context.Logger().Print(fmt.Errorf("Cannot add new API Token in Cache, warning => %v", err))
		//QUESTION: possible to put consistency flyweight here?
	}
	err = addTokenToCacheList(ID, token)
	if err != nil {
		context.Logger().Print(fmt.Errorf("Cannot add new API Token in Developer list in Cache, warning => %v", err))
		//QUESTION: possible to put consistency flyweight here?
	}
	responseFromServer := developerResponses.AddToken{}
	responseFromServer.FromAPIToken(token)
	return context.JSON(http.StatusCreated, &responseFromServer)
}

//HandleDropAPI_Token handles a request to remove a developer API Token.
func HandleDropAPI_Token(context echo.Context) error {
	request := new(developerRequests.DropToken)
	err := request.FromForm(context)
	if err != nil {
		errorResp := errorResponses.ErrorDetail{}
		errorResp.FromError(err, http.StatusBadRequest)
		return context.JSON(http.StatusBadRequest, &errorResp)
	}
	context.Logger().Print(request)
	IsValid, err := controllerSharedFuncs.IsValidAPI_Token(request.API_Token)
	if !IsValid || err != nil {
		context.Logger().Print(fmt.Errorf("API Token %s rejected", request.API_Token))
		errorResp := errorResponses.ErrorDetail{}
		errorResp.FromError(errors.New("Rejected by the system"), http.StatusBadRequest)
		return context.JSON(http.StatusBadRequest, &errorResp)
	}

	ID, err := getDevIDFromSessionToken(request.SessionToken)
	if err != nil {
		errorResp := errorResponses.ErrorDetail{}
		context.Logger().Print(fmt.Errorf("%s token rejected by the system, Invalid Session", request.SessionToken))
		errorResp.FromError(errors.New("Rejected by the system"), http.StatusBadRequest)
		return context.JSON(http.StatusBadRequest, &errorResp)
	}

	cacheCleared, err := removeAPI_Token(ID, request.TokenToDrop)
	if err != nil {
		if cacheCleared {
			//just Log and return error
			context.Logger().Print("Cache error : see below")
		} else {
			//more like a warning : cache is ok but archives are not. which is not consistent.
			context.Logger().Print("Warning from Archives : see below")
		}
		context.Logger().Print(fmt.Errorf("%s API Token not removed. Error => %v", request.TokenToDrop, err))
		errorResp := errorResponses.ErrorDetail{}
		errorResp.FromError(errors.New("Cannot remove API Token"), http.StatusInternalServerError)
		return context.JSON(http.StatusInternalServerError, errorResp)
	}

	responseFromServer := developerResponses.DropToken{}
	responseFromServer.FromOldAPIToken(request.TokenToDrop)
	return context.JSON(http.StatusOK, &responseFromServer)
}

//HandleRegistration handles a request to register a developer.
func HandleRegistration(context echo.Context) error {
	request := new(developerRequests.DevRegistration)
	err := request.FromForm(context)
	if err != nil {
		errorResp := errorResponses.ErrorDetail{}
		errorResp.FromError(err, http.StatusBadRequest)
		//errorResp.ErrorMessage += fmt.Sprintf("%v", context)
		return context.JSON(http.StatusBadRequest, &errorResp)
	}

	IsValid, err := controllerSharedFuncs.IsValidAPI_Token(request.API_Token)
	if !IsValid || err != nil {
		context.Logger().Print(fmt.Errorf("API Token %s rejected : error %v", request.API_Token, err))
		errorResp := errorResponses.ErrorDetail{}
		errorResp.FromError(errors.New("Rejected by the system"), http.StatusBadRequest)
		return context.JSON(http.StatusBadRequest, &errorResp)
	}

	ID, err := registerDeveloper(*request)
	if err != nil {
		errorResp := errorResponses.ErrorDetail{}
		errorResp.FromError(err, http.StatusInternalServerError)
		return context.JSON(http.StatusBadRequest, &errorResp)
	}

	token, err := updateCacheWithSessionDeveloperToken(ID)
	if err != nil {
		errorResp := errorResponses.ErrorDetail{}
		errorResp.FromError(errors.New("User registered, but I did not login automatically, try to login later"), http.StatusBadRequest)
		return context.JSON(http.StatusInternalServerError, errorResp)
	}

	responseFromServer := developerResponses.DevAuth{}
	responseFromServer.FromToken(token)
	return context.JSON(http.StatusCreated, &responseFromServer)
}

//HandleLogin handles login requests for developers.
func HandleLogin(context echo.Context) error {
	request := new(developerRequests.DevAuth)
	err := request.FromForm(context)
	if err != nil {
		errorResp := errorResponses.ErrorDetail{}
		errorResp.FromError(err, http.StatusBadRequest)
		return context.JSON(http.StatusBadRequest, &errorResp)
	}

	IsValid, err := controllerSharedFuncs.IsValidAPI_Token(request.API_Token)
	if !IsValid || err != nil {
		context.Logger().Print(fmt.Errorf("API Token rejected %v : error %v", request.API_Token, err))
		errorResp := errorResponses.ErrorDetail{}
		errorResp.FromError(errors.New("Rejected by the system"), http.StatusBadRequest)
		return context.JSON(http.StatusBadRequest, &errorResp)
	}

	isLoggable, developerID, err := checkLogin(*request)
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

	token, err := updateCacheWithSessionDeveloperToken(developerID)
	if err != nil {
		errorResp := errorResponses.ErrorDetail{}
		context.Logger().Print(err)
		errorResp.FromError(errors.New("Temporary error, retry in a few seconds"), http.StatusInternalServerError)
		return context.JSON(http.StatusInternalServerError, errorResp)
	}
	responseFromServer := developerResponses.DevAuth{}
	responseFromServer.FromToken(token)
	return context.JSON(http.StatusCreated, &responseFromServer)
}
