package OutGameRequests

import (
	"errors"

	"github.com/labstack/echo"
)

//GameList represents a request to all the games available to a single user.
type GameList struct {
	Type         string `json:"Type" xml:"Type" form:"Type"`
	API_Token    string `json:"API_Token" xml:"API_Token" form:"API_Token"`
	SessionToken string `json:"SessionToken" xml:"SessionToken" form:"SessionToken"`
}

//FromForm creates a valid Sruct based on form data submitted, or returns error.
//
// Does not check for the validity of the items inside the struct (e.g. tokens)
func (receiver *GameList) FromForm(c echo.Context) error {
	var err error
	receiver.Type = c.FormValue("Type")
	receiver.API_Token = c.FormValue("API_Token")
	receiver.SessionToken = c.FormValue("SessionToken")
	if receiver.Type != "GameList" || receiver.API_Token == "" || receiver.SessionToken == "" {
		err = errors.New("Invalid Form Submitted")
	}

	return err
}
