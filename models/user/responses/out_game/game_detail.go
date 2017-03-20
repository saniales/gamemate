package outGameResponses

import (
	"sanino/gamemate/models/user/data_structures"
)

//UserGameDetail represents a request to get details for a single game.
type UserGameDetail struct {
	Type string                         `json:"Type" xml:"Type" form:"Type"`
	Game userDataStructs.SummarizedGame `json:"Game" xml:"Game" form:"Game"`
}

//FromGame creates a valid Struct based on the passed game, or returns error.
//
//Does not check for the validity of the items inside the struct (e.g. tokens)
func (receiver *UserGameDetail) FromGame(Game userDataStructs.SummarizedGame) error {
	receiver.Type = "UserGameDetail"
	receiver.Game = Game
	return nil
}
