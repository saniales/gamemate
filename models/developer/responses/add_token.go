package developerResponses

import "strings"

//AddToken represents a POSITIVE response from the server to a developerRequests.AddToken
//
//For NEGATIVE response, please refer to errorResponses.ErrorResponse.
type AddToken struct {
	Type         string `json:"Type" xml:"Type" form:"Type"`
	NewAPI_Token string `json:"NewAPI_Token" xml:"NewAPI_Token" form:"NewAPI_Token"`
}

//FromAPIToken fills the struct's data with proper definition, based on an
//API token.
func (receiver *AddToken) FromAPIToken(API_Token string) {
	receiver.Type = "AddToken"
	receiver.NewAPI_Token = strings.ToUpper(strings.Replace(API_Token, "0x", "", 1))
}
