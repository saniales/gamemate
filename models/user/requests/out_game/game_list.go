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

func (receiver *GameList) FromForm(c echo.Context) error {
	var err error
	receiver.Type = c.FormValue("Type")
	receiver.API_Token = c.FormValue("API_Token")
	receiver.SessionToken = c.FormValue("SessionToken")
	if receiver.Type == "" || receiver.API_Token == "" || receiver.SessionToken == "" {
		err = errors.New("Invalid Form Submitted, cannot find some fields")
	}

	return err
}
