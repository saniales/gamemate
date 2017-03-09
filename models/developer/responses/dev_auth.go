package developerResponses

//DevAuth represents the POSITIVE response to a registration or login request.
//
//For NEGATIVE response, please refer to errorResponses.ErrorResponse.
type DevAuth struct {
	Type         string `json:"Type" xml:"Type" form:"Type"`
	SessionToken string `json:"SessionToken" xml:"SessionToken" form:"SessionToken"`
}

//FromToken fills the struct's data with proper definition, based on a session token.
func (receiver *DevAuth) FromToken(SessionToken string) {
	receiver.Type = "DevSessionToken"
	receiver.SessionToken = SessionToken
}
