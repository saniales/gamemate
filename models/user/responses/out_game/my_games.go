package outGameResponses

import (
	"errors"
	"sanino/gamemate/models/user/data_structures"
)

//MyGames represents a list of games available to an user, which made a SUCCESSFULL request.
//
//For NEGATIVE response, please refer to errorResponses.ErrorResponse.
type MyGames struct {
	Type  string                           `json:"Type" xml:"Type" form:"Type"`
	Games []userDataStructs.SummarizedGame `json:"Games" xml:"Games" form:"Games"`
}

//FromGames creates a GameList from a list of games.
func (receiver *MyGames) FromGames(Games []userDataStructs.SummarizedGame) error {
	receiver.Type = "Game List"
	if Games == nil {
		return errors.New("Assigning a nil set of games; you can't do it")
	}
	receiver.Games = Games
	return nil
}

//Game represents a game saved into the system.
type Game struct {
	Type       string
	ID         int64
	Name       string
	NumPlayers int64
}
