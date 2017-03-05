package developerRequests

import (
	"errors"

	"github.com/labstack/echo"
)

//DevRegistration represents a request to register a deveoper into the system.
type DevRegistration struct {
	Type      string `json:"Type" xml:"Type" form:"Type"`
	API_Token string `json:"API_Token" xml:"API_Token" form:"API_Token"`
	Email     string `json:"Email" xml:"Email" form:"Email"`
	Password  string `json:"Password" xml:"Password" form:"Password"`
}

// FromForm Converts from a submitted form (or request) to his struct.
//
// Does not check for the validity of the items inside the struct (e.g. tokens)
func (receiver *DevRegistration) FromForm(c echo.Context) error {
	var err error
	/*
		receiver.Type = c.FormValue("Type")
		receiver.API_Token = c.FormValue("API_Token")
		receiver.Password = c.FormValue("Password")
		receiver.Email = c.FormValue("Email")
	*/
	err = c.Bind(receiver)
	if receiver.Type != "DevRegistration" || err != nil {
		err = errors.New("Invalid Form Submitted " + err.Error())
	}
	return err
}
