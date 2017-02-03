package libs


//Represents an auth try to the system.
type Auth struct {
    Username string `json:"Username" xml:"Username" form:"Username"`
    Password string `json:"Password" xml:"Password" form:"Password"`
}
