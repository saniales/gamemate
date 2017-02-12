package vendorResponses

//AddGame represents a POSITIVE response to a vendorRequests.AddGame
//
//For NEGATIVE response, please refer to errorResponses.ErrorResponse.
type AddGame struct {
	Type   string `json:"Type" xml:"Type" form:"Type"`
	GameID int64  `json:"GameID" xml:"GameID" form:"GameID"`
	Result string `json:"Result" xml:"Result" form:"Result"`
}

//FromGameID fills the struct with data from a game ID.
func (receiver *AddGame) FromGameID(GameID int64) {
	receiver.Type = "AddGame"
	receiver.GameID = GameID
	receiver.Result = "OK"
}
