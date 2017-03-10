package developerResponses

//TokenList represents the POSITIVE response to a token_list request.
//
//For NEGATIVE response, please refer to errorResponses.ErrorResponse.
type TokenList struct {
	Type   string   `json:"Type" xml:"Type" form:"Type"`
	Tokens []string `json:"Tokens" xml:"Tokens" form:"Tokens"`
}

//FromTokens fills the struct's data with proper definition, based on a provided API tokens.
func (receiver *TokenList) FromTokens(Tokens []string) {
	receiver.Type = "DevTokenList"
	receiver.Tokens = Tokens
}
