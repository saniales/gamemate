package outGameResponses

//GameList represents a list of games available to an user, which made a SUCCESSFULL request.
//
//For NEGATIVE response, please refer to errorResponses.ErrorResponse.
type GameList struct {
	Type  string
	Games []Game
}

//FromGames creates a GameList from a list of games.
func (receiver *GameList) FromGames(Games []Game) {
	receiver.Type = "Game List"
	receiver.Games = Games
}

//Game represents a game saved into the system.
type Game struct {
	Type       string
	ID         int64
	Name       string
	NumPlayers int64
}
