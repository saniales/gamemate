package vendorRequests

import (
	"errors"
	"strconv"

	"github.com/labstack/echo"
)

//RemoveGame represents a request to remove a game from a vendors list.
//
//Only the owner of the game can perform this action.
type RemoveGame struct {
	Type         string `json:"Type" xml:"Type" form:"Type"`
	API_Token    string `json:"API_Token" xml:"API_Token" form:"API_Token"`
	SessionToken string `json:"SessionToken" xml:"SessionToken" form:"SessionToken"`
	GameID       int64  `json:"GameID" xml:"GameID" form:"GameID"`
}

//FromForm creates a valid Struct based on form data submitted, or returns error.
//
// Does not check for the validity of the items inside the struct (e.g. tokens)
func (receiver *RemoveGame) FromForm(c echo.Context) error {
	var err error
	receiver.Type = c.FormValue("Type")
	receiver.API_Token = c.FormValue("API_Token")
	receiver.SessionToken = c.FormValue("SessionToken")
	receiver.GameID, err = strconv.ParseInt(c.FormValue("GameID"), 10, 64)
	if err != nil {
		return errors.New("Invalid Form Submitted")
	}

	if receiver.Type != "RemoveGame" || receiver.API_Token == "" ||
		receiver.SessionToken == "" || receiver.GameID <= 0 {
		return errors.New("Invalid Form Submitted")
	}

	return err
}
