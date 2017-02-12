package loginRequests

import (
	"errors"

	"github.com/labstack/echo"
)

//Auth represents an auth try to the system.
type Auth struct {
	Type      string `json:"Type" xml:"Type" form:"Type"`
	API_Token string `json:"API_Token" xml:"API_Token" form:"API_Token"`
	Username  string `json:"Username" xml:"Username" form:"Username"`
	Password  string `json:"Password" xml:"Username" form:"Password"`
}

// FromForm Converts from a submitted form (or request) to his struct.
//
// Does not check for the validity of the items inside the struct (e.g. tokens)
func (receiver *Auth) FromForm(c echo.Context) error {
	var err error
	receiver.Type = c.FormValue("Type")
	receiver.API_Token = c.FormValue("API_Token")
	receiver.Username = c.FormValue("Username")
	receiver.Password = c.FormValue("Password")

	if receiver.Type != "Auth" || receiver.Username == "" || receiver.Password == "" || receiver.API_Token == "" {
		err = errors.New("Invalid Form Submitted")
	}
	return err
}
