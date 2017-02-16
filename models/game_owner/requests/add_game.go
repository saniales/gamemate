package gameOwnerRequests

import (
	"errors"

	"github.com/labstack/echo"
)

//AddGame represents a request to add a token for an app of the developer.
type AddGame struct {
	Type            string `json:"Type" xml:"Type" form:"Type"`
	API_Token       string `json:"API_Token" xml:"API_Token" form:"API_Token"`
	SessionToken    string `json:"SessionToken" xml:"SessionToken" form:"SessionToken"`
	GameName        string `json:"GameName" xml:"GameName" form:"GameName"`
	GameDescription string `json:"GameDescription" xml:"GameDescription" form:"GameDescription"`
	MatchMaxPlayers int64  `json:"MatchMaxPlayers" xml:"MatchMaxPlayers" form:"MatchMaxPlayers"`
}

//FromForm creates a valid Struct based on form data submitted, or returns error.
//
// Does not check for the validity of the items inside the struct (e.g. tokens)
func (receiver *AddGame) FromForm(c echo.Context) error {
	receiver.Type = c.FormValue("Type")
	receiver.API_Token = c.FormValue("API_Token")
	receiver.SessionToken = c.FormValue("SessionToken")
	receiver.GameName = c.FormValue("GameName")
	receiver.GameDescription = c.FormValue("GameDescription")
	if receiver.Type != "AddGame" || receiver.API_Token == "" ||
		receiver.SessionToken == "" || receiver.GameName == "" ||
		receiver.GameDescription == "" {
		return errors.New("Invalid Form Submitted")
	}

	return nil
}
