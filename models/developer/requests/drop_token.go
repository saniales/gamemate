package developerRequests

import (
	"errors"

	"github.com/labstack/echo"
)

//DropToken represents a POSITIVE response to a developerRequests.DropToken.
//
//If the response is NEGATIVE, please refer to error
type DropToken struct {
	Type         string `json:"Type" xml:"Type" form:"Type"`
	API_Token    string `json:"API_Token" xml:"API_Token" form:"API_Token"`
	SessionToken string `json:"SessionToken" xml:"SessionToken" form:"SessionToken"`
	TokenToDrop  string `json:"TokenToDrop" xml:"TokenToDrop" form:"TokenToDrop"`
}

// FromForm Converts from a submitted form (or request) to his struct.
//
// Does not check for the validity of the items inside the struct (e.g. tokens)
func (receiver *DropToken) FromForm(c echo.Context) error {
	var err error
	receiver.Type = c.FormValue("Type")
	receiver.API_Token = c.FormValue("API_Token")
	receiver.SessionToken = c.FormValue("SessionToken")
	receiver.TokenToDrop = c.FormValue("TokenToDrop")

	if receiver.Type != "DropToken" || receiver.SessionToken == "" ||
		receiver.TokenToDrop == "" || receiver.API_Token == "" {
		err = errors.New("Invalid Form Submitted")
	}
	return err
}
