package vendorResponses

import (
	"data_structures"
	"errors"
)

//MyGames represents a response from the server to a vendorResponses.MyGames
type MyGames struct {
	Type  string                   `json:"Type" xml:"Type" form:"Type"`
	Games []vendorDataStructs.Game `json:"Games" xml:"Games" form:"Games"`
}

//FromGames fills the structs data from a list of Games.
func (receiver *MyGames) FromGames(Games []vendorDataStructs.Game) error {
	receiver.Type = "MyGamesVendor"
	receiver.Games = Games
	if receiver.Games == nil {
		return errors.New("FromGames() error : Passed nil argument Games")
	}
	return nil
}
