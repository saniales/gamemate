package loginRequests

import (
	"errors"
	"time"

	"github.com/labstack/echo"
)

//Registration represents a request to register into the system.
type Registration struct {
	Type      string `json:"Type" xml:"Type" form:"Type"`
	API_Token string `json:"API_Token" xml:"API_Token" form:"API_Token"`
	Username  string `json:"Username" xml:"Username" form:"Username"`
	Email     string `json:"Email" xml:"Email" form:"Email"`
	Password  string `json:"Password" xml:"Password" form:"Password"`
	Birthday  string `json:"Birthday" xml:"Birthday" form:"Birthday"`
	Gender    string `json:"Gender" xml:"Gender" form:"Gender"`
}

// FromForm Converts from a submitted form (or request) to his struct.
//
// Does not check for the validity of the items inside the struct (e.g. tokens)
func (receiver *Registration) FromForm(c echo.Context) error {
	err := c.Bind(receiver)
	if err != nil || receiver.Type != "Registration" {
		return errors.New("Invalid Form Submitted " + err.Error())
	}
	_, err = receiver.BirthdayDate()
	if err != nil {
		return errors.New("Invalid Form Submitted, Birthday is not in a correct format => " + receiver.Birthday)
	}
	return nil
}

//BirthdayDate converts the date string in a time struct
func (receiver *Registration) BirthdayDate() (time.Time, error) {
	return time.Parse("2006-01-02", receiver.Birthday)
}
