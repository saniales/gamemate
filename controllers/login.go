package controllers

import (
	"sanino/gamemate/models"

	"github.com/labstack/echo"
)

//HandleAuth handles the authentication of the user for the system.
func HandleAuth(c echo.Context) error {
	var AuthTry = models.Auth{}
	var err = AuthTry.FromForm(c)

	if err != nil {
		return err
	}
	// answers with session token valid for 30 minutes from last session
	// it has to be created and put in redis after correct authentication.
	// otherwise the system must reply with a system_error struct {Code : 1}
	// if debug it includes a message {errCode : 1, message : "ZIOBANANA"}
	return nil
}

//HandleRegistration handles the registration of a user for the system.
func HandleRegistration(context echo.Context) error {
	var RegTry = models.Registration{}
	var err = RegTry.FromForm(context)
	if err != nil {
		return err
	}
	//else query and if query successful add user also into cache and reply with session_token

	return nil
}

//HandlePlayerInfoRequest handles a player request of information about a player.
func HandlePlayerInfoRequest(context echo.Context) error {
	return nil
}
