package outGameResponses

//MyEnabledGames represents game IDs for games available to the user to play.
type MyEnabledGames struct {
	Type  string  `json:"Type" xml:"Type" form:"Type"`
	Games []int64 `json:"Games" xml:"Games" form:"Games"`
}

//FromGameIDs fills the struct with the gameIDs data.
func (receiver *MyEnabledGames) FromGameIDs(gameIDs []int64) {
	receiver.Type = "My Enabled Games"
	receiver.Games = gameIDs
}
