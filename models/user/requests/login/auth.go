package loginRequests

import (
	"errors"

	"github.com/labstack/echo"
)

//UserAuth represents an auth try to the system.
type UserAuth struct {
	Type      string `json:"Type" xml:"Type" form:"Type"`
	API_Token string `json:"API_Token" xml:"API_Token" form:"API_Token"`
	Username  string `json:"Username" xml:"Username" form:"Username"`
	Password  string `json:"Password" xml:"Username" form:"Password"`
}

// FromForm Converts from a submitted form (or request) to his struct.
//
// Does not check for the validity of the items inside the struct (e.g. tokens)
func (receiver *UserAuth) FromForm(c echo.Context) error {
	err := c.Bind(receiver)
	if err != nil || receiver.Type != "UserAuth" {
		return errors.New("Invalid Form Submitted " + err.Error())
	}
	return nil
}
