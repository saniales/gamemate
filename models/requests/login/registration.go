package loginRequests

import (
	"errors"
	"time"

	"github.com/labstack/echo"
)

//Registration represents a request to register into the system.
type Registration struct {
	Username string `json:"Username" xml:"Username" form:"Username"`
	Email    string `json:"Email" xml:"Email" form:"Email"`
	Password string `json:"Password" xml:"Password" form:"Password"`
	Birthday string `json:"Birthday" xml:"Birthday" form:"Birthday"`
	Gender   string `json:"Gender" xml:"Gender" form:"Gender"`
}

// FromForm Converts from a submitted form (or request) to his struct.
func (receiver *Registration) FromForm(c echo.Context) error {
	var err error
	receiver.Username = c.FormValue("Username")
	receiver.Password = c.FormValue("Password")
	receiver.Email = c.FormValue("Email")
	receiver.Birthday = c.FormValue("Birthday")
	receiver.Gender = c.FormValue("Gender")
	if receiver.Username == "" || receiver.Password == "" || receiver.Email == "" || receiver.Birthday == "" || receiver.Gender == "" {
		err = errors.New("Invalid Form Submitted, cannot find some fields")
	} else if _, err = receiver.BirthdayDate(); err != nil {
		err = errors.New("Invalid Form Submitted, Birthday is not in a correct format => " + receiver.Birthday)
	}
	return err
}

//BirthdayDate converts the date string in a time struct
func (receiver *Registration) BirthdayDate() (time.Time, error) {
	return time.Parse("2006-01-02", receiver.Birthday)
}
