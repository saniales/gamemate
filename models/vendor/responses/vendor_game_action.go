package vendorResponses

//VendorGameAction represents a POSITIVE response to a vendorRequests.VendorGameAction
//
//For NEGATIVE response, please refer to errorResponses.ErrorResponse.
type VendorGameAction struct {
	Type   string `json:"Type" xml:"Type" form:"Type"`
	GameID int64  `json:"GameID" xml:"GameID" form:"GameID"`
	Result string `json:"Result" xml:"Result" form:"Result"`
}

//FromGameID fills the struct with a POSITIVE response from a GameID.
func (receiver *VendorGameAction) FromGameID(GameID int64) {
	receiver.Type = "VendorGameAction"
	receiver.GameID = GameID
	receiver.Result = "OK"
}
