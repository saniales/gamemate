package models

import (
	"errors"

	"github.com/labstack/echo"
)

//Auth represents an auth try to the system.
type Auth struct {
	Username string `json:"Username" xml:"Username" form:"Username"`
	Password string `json:"Password" xml:"Username" form:"Password"`
}

// FromForm Converts from a submitted form (or request) to his struct.
func (receiver *Auth) FromForm(c echo.Context) error {
	var err error
	receiver.Username = c.FormValue("Username")
	receiver.Password = c.FormValue("Password")
	if receiver.Username == "" || receiver.Password == "" {
		err = errors.New("Invalid Form Submitted, cannot find username or password.")
	}
	return err
}
