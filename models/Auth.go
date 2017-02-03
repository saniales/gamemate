package models

import (
    "github.com/labstack/echo"
)

//Represents an auth try to the system.
type Auth struct {
    Username string `json:"Username" xml:"Username" form:"Username"`
    Password string `json:"Password" xml:"Username" form:"Password"`
}

func (this *Auth) FromForm(c echo.Context) (FormDecodable, error) {
    var err error = nil
    this.Username = c.FormValue("Username")
    this.Password = c.FormValue("Password")
    if this.Username == "" || this.Password == "" {
        err = errors.New("Invalid Form Submitted, cannot find username or password.")
    }
    return this, err
}
