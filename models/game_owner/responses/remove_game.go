package gameOwnerResponses

//RemoveGame represents a response to a gameOwnerRequests.RemoveGame.
//
//Only the owner of the game can perform this action.
//For NEGATIVE response, please refer to errorResponses.ErrorResponse.
type RemoveGame struct {
	Type   string `json:"Type" xml:"Type" form:"Type"`
	GameID int64  `json:"GameID" xml:"GameID" form:"GameID"`
	Result string `json:"Result" xml:"Result" form:"Result"`
}

//FromGameID fills the struct with a POSITIVE response from a GameID.
func (receiver *RemoveGame) FromGameID(GameID int64) {
	receiver.Type = "GameOwnerGameAction"
	receiver.GameID = GameID
	receiver.Result = "OK"
}
