package outGameRequests

import (
	"errors"

	"github.com/labstack/echo"
)

//UserGameDetail represents a request to get details for a single game.
type UserGameDetail struct {
	Type         string `json:"Type" xml:"Type" form:"Type"`
	API_Token    string `json:"API_Token" xml:"API_Token" form:"API_Token"`
	SessionToken string `json:"SessionToken" xml:"SessionToken" form:"SessionToken"`
	GameID       int64  `json:"GameID" xml:"GameID" form:"GameID"`
}

//FromForm creates a valid Sruct based on form data submitted, or returns error.
//
// Does not check for the validity of the items inside the struct (e.g. tokens)
func (receiver *UserGameDetail) FromForm(c echo.Context) error {
	err := c.Bind(receiver)
	if err != nil || receiver.Type != "UserGameDetail" {
		return errors.New("Invalid Form Submitted")
	}
	return nil
}
