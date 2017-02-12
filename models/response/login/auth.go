package response

//Auth represents the response to a registration or login request.
type Auth struct {
	SessionToken string `json:"SessionToken" xml:"SessionToken" form:"SessionToken"`
}
