package request

import (
	"errors"
	"time"

	"github.com/labstack/echo"
)

//Registration represents a request to register into the system.
type Registration struct {
	Username string    `json:"Username" xml:"Username" form:"Username"`
	Email    string    `json:"Email" xml:"Email" form:"Email"`
	Password string    `json:"Password" xml:"Password" form:"Password"`
	Birthday time.Time `json:"Birthday" xml:"Birthday" form:"Birthday"`
	Gender   string    `json:"Gender" xml:"Gender" form:"Gender"`
}

// FromForm Converts from a submitted form (or request) to his struct.
func (receiver *Registration) FromForm(c echo.Context) error {
	var err, errBDay error
	receiver.Username = c.FormValue("Username")
	receiver.Password = c.FormValue("Password")
	receiver.Email = c.FormValue("Email")
	receiver.Birthday, errBDay = time.Parse(time.RFC3339, c.FormValue("Birthday"))
	receiver.Gender = c.FormValue("Gender")
	if receiver.Username == "" || receiver.Password == "" || receiver.Email == "" || errBDay != nil || receiver.Gender == "" {
		err = errors.New("Invalid Form Submitted, cannot find some fields")
	}
	return err
}
