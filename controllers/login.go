package controllers

import(
  "github.com/labstack/echo"
  "Tesi/constants"
)


func HandleAuth(context echo.Context) error {
    var err error
    var AuthTry Auth
    AuthTry, err = new(Auth).FromForm(c)
    if err != nil {
        // answers with session token valid for 30 minutes from last session
        // it has to be created and put in redis after correct authentication.
        // otherwise the system must reply with a system_error struct {errCode : 1}
        // if debug it includes a message {errCode : 1, message : "ZIOBANANA"}
        
    } else {

    }

}

func HandleRegistration(context echo.Context) error {

}

func HandlePlayerInfoRequest(context echo.Context) error {

}
