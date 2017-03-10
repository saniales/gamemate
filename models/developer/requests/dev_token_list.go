package developerRequests

import (
	"errors"

	"github.com/labstack/echo"
)

//TokenList represents a request to add a token for an app of the developer.
type TokenList struct {
	Type         string `json:"Type" xml:"Type" form:"Type"`
	API_Token    string `json:"API_Token" xml:"API_Token" form:"API_Token"`
	SessionToken string `json:"SessionToken" xml:"SessionToken" form:"SessionToken"`
}

//FromForm creates a valid Struct based on form data submitted, or returns error.
//
// Does not check for the validity of the items inside the struct (e.g. tokens)
func (receiver *TokenList) FromForm(c echo.Context) error {
	err := c.Bind(receiver)
	if err != nil || receiver.Type != "DevTokenList" {
		err = errors.New("Invalid Form Submitted")
	}

	return err
}
