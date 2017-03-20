package outGameResponses

import (
	"errors"
	"sanino/gamemate/models/user/data_structures"
)

//UserGameList represents a list of games available to an user, which made a SUCCESSFULL request.
//
//For NEGATIVE response, please refer to errorResponses.ErrorResponse.
type UserGameList struct {
	Type  string                           `json:"Type" xml:"Type" form:"Type"`
	Games []userDataStructs.SummarizedGame `json:"Games" xml:"Games" form:"Games"`
}

//FromGames creates a GameList from a list of games.
func (receiver *UserGameList) FromGames(Games []userDataStructs.SummarizedGame) error {
	receiver.Type = "UserGameList"
	if Games == nil {
		return errors.New("Assigning a nil set of games")
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
