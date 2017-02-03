package libs

//Represents a request to register into the system.
type Registration struct {
    Username string `json:"Username" xml:"Username" form:"Username"`
    Email    string `json:"Email" xml:"Email" form:"Email"`
    Password string `json:"Password" xml:"Password" form:"Password"`
    Birthday int64  `json:"Birthday" xml:"Birthday" form:"Birthday"`
    Gender   string `json:"Gender" xml:"Gender" form:"Gender"`
}
