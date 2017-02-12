package developerResponses

//AddToken represents a POSITIVE response from the server to a developerRequests.AddToken
//
//For NEGATIVE response, please refer to errorResponses.ErrorResponse.
type AddToken struct {
	NewAPI_Token string `json:"SessionToken" xml:"SessionToken" form:"SessionToken"`
}

//FromAPIToken fills the struct's data with proper definition, based on an
//API token.
func (receiver *AddToken) FromAPIToken(API_Token string) {
	receiver.Type = "Developer New API Token"
	receiver.NewAPI_Token = API_Token
}
