package gameServerLogic

type Result int

const (
	WIN     Result = iota
	DRAW    Result = iota
	ONGOING Result = iota
)

//MoveChecker Represents an entity which controls validity of moves for a game,
//and apply those moves.
type MoveChecker interface {
	IsValidMove(map[string]interface{}) bool //IsValidMove checks if a move is valid.
	MakeMove(map[string]interface{}) bool    //MakeMove does a move, if it is valid.
	CheckWin(map[string]interface{}) Result  //CheckWin checks if the last move won the game, from the last move.
}
