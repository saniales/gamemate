package sharedResponses

//GameOwnerGameAction represents a POSITIVE response to a gameOwnerRequests.GameOwnerGameAction
//
//For NEGATIVE response, please refer to errorResponses.ErrorResponse.
type GameAction struct {
	Type   string `json:"Type" xml:"Type" form:"Type"`
	GameID int64  `json:"GameID" xml:"GameID" form:"GameID"`
	Result string `json:"Result" xml:"Result" form:"Result"`
}

//FromGameID fills the struct with a POSITIVE response from a GameID.
func (receiver *GameAction) FromGameID(GameID int64) {
	receiver.Type = "GameAction"
	receiver.GameID = GameID
	receiver.Result = "OK"
}
