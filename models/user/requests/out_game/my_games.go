package outGameRequests

import (
	"errors"

	"github.com/labstack/echo"
)

//MyGames represents a request dome (or all) the games available to a single user.
type MyGames struct {
	Type         string `json:"Type" xml:"Type" form:"Type"`
	API_Token    string `json:"API_Token" xml:"API_Token" form:"API_Token"`
	SessionToken string `json:"SessionToken" xml:"SessionToken" form:"SessionToken"`
}

//FromForm creates a valid Sruct based on form data submitted, or returns error.
//
// Does not check for the validity of the items inside the struct (e.g. tokens)
func (receiver *MyGames) FromForm(c echo.Context) error {
	err := c.Bind(receiver)
	if err != nil || receiver.Type != "GameList" {
		return errors.New("Invalid Form Submitted " + err.Error())
	}
	return nil
}
