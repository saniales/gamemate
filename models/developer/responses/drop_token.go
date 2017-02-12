package developerResponses

//DropToken represents a POSITIVE response from the server to a
//developerRequests.DropToken.
//
//For NEGATIVE response, please refer to errorResponses.ErrorResponse.
type DropToken struct {
	Type     string `json:"Type" xml:"Type" form:"Type"`
	OldToken string `json:"OldToken" xml:"OldToken" form:"OldToken"`
	Result   string `json:"Result" xml:"Result" form:"Result"`
}

//FromOldAPIToken fills the struct's data with proper definition, based on an old
//API token.
func (receiver *DropToken) FromOldAPIToken(API_Token string) {
	receiver.Type = "Developer New API Token"
	receiver.OldToken = API_Token
	receiver.Result = "OK"
}
