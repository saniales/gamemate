package gameOwnerRequests

import (
	"errors"

	"github.com/labstack/echo"
)

//RemoveGame represents a request to remove a game from a gameOwners list.
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
	err := c.Bind(receiver)
	if err != nil || receiver.Type != "RemoveGame" {
		return errors.New("Invalid Form Submitted " + err.Error())
	}

	return err
}
