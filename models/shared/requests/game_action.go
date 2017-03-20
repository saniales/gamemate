package sharedRequests

import (
	"errors"

	"github.com/labstack/echo"
)

//GameAction represents a request to enable a game for a user.
type GameAction struct {
	Type         string `json:"Type" xml:"Type" form:"Type"`
	API_Token    string `json:"API_Token" xml:"API_Token" form:"API_Token"`
	SessionToken string `json:"SessionToken" xml:"SessionToken" form:"SessionToken"`
	GameID       int64  `json:"GameID" xml:"GameID" form:"GameID"`
	UserID       int64  `json:"UserID" xml:"UserID" form:"UserID"`
	Action       bool   `json:"Action" xml:"Action" form:"Action"` //True = Enable, False = Disable (on a GameID)
}

//FromForm creates a valid Sruct based on form data submitted, or returns error.
//
// Does not check for the validity of the items inside the struct (e.g. tokens)
func (receiver *GameAction) FromForm(c echo.Context) error {
	err := c.Bind(receiver)
	if err != nil || receiver.Type != "GameAction" || receiver.GameID <= 0 ||
		receiver.UserID < -1 { //userID = -1 means "requesting user"
		return errors.New("Invalid Form Submitted")
	}

	return nil
}
