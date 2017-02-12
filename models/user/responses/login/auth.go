package loginResponses

//Auth represents the response to a registration or login request.
type Auth struct {
	Type         string `json:"Type" xml:"Type" form:"Type"`
	SessionToken string `json:"SessionToken" xml:"SessionToken" form:"SessionToken"`
}

//FromToken fills the struct's data with proper definition, based on a session token.
func (receiver *Auth) FromToken(SessionToken string) {
	receiver.Type = "Session Token"
	receiver.SessionToken = SessionToken
}
