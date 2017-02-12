package vendorRequests

import (
	"errors"
	"strconv"

	"github.com/labstack/echo"
)

//VendorGameAction represents a request to enable a game for a user.
type VendorGameAction struct {
	Type         string `json:"Type" xml:"Type" form:"Type"`
	API_Token    string `json:"API_Token" xml:"API_Token" form:"API_Token"`
	SessionToken string `json:"SessionToken" xml:"SessionToken" form:"SessionToken"`
	GameID       int64  `json:"GameID" xml:"GameID" form:"GameID"`
	UserID       int64  `json:"UserID" xml:"UserID" form:"UserID"`
	Action       bool   `json:"Action" xml:"Action" form:"Action"` //True = Enable, False = Disable (on a GameID)
}

//FromForm creates a valid Sruct based on form data submitted, or returns error.
//
// Does not check for the validity of the items inside the struct (e.g. tokens)
func (receiver *VendorGameAction) FromForm(c echo.Context) error {
	var err error
	errMsg := "Invalid Form Submitted"
	receiver.Type = c.FormValue("Type")
	receiver.API_Token = c.FormValue("API_Token")
	receiver.SessionToken = c.FormValue("SessionToken")

	receiver.GameID, err = strconv.ParseInt(c.FormValue("GameID"), 10, 64)
	if err != nil {
		return errors.New(errMsg)
	}

	receiver.UserID, err = strconv.ParseInt(c.FormValue("UserID"), 10, 64)
	if err != nil {
		return errors.New(errMsg)
	}

	receiver.Action, err = strconv.ParseBool(c.FormValue("Action"))
	if err != nil {
		return errors.New(errMsg)
	}

	if receiver.Type != "VendorGameAction" || receiver.API_Token == "" ||
		receiver.SessionToken == "" || receiver.GameID <= 0 ||
		receiver.UserID <= 0 {
		return errors.New(errMsg)
	}

	return err
}
