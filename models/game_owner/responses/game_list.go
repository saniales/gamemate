package gameOwnerResponses

import (
	"errors"
	"sanino/gamemate/models/game_owner/data_structures"
)

//MyGames represents a response from the server to a gameOwnerRequests.MyGames
type GameOwnerGameList struct {
	Type  string                      `json:"Type" xml:"Type" form:"Type"`
	Games []gameOwnerDataStructs.Game `json:"Games" xml:"Games" form:"Games"`
}

//FromGames fills the structs data from a list of Games.
func (receiver *GameOwnerGameList) FromGames(Games []gameOwnerDataStructs.Game) error {
	receiver.Type = "GameOwnerGameList"
	receiver.Games = Games
	if receiver.Games == nil {
		return errors.New("FromGames() error : Passed nil argument Games")
	}
	return nil
}
