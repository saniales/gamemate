package gameServerLogic

type TicTacToeSymbol int

const (
	EMPTY_CELL  TicTacToeSymbol = iota
	CROSS_CELL  TicTacToeSymbol = iota
	CIRCLE_CELL TicTacToeSymbol = iota
)

//TicTacToeChecker Represents the MoveChecker for Tic Tac Toe game.
type TicTacToeChecker struct {
	gameGrid [3][3]TicTacToeSymbol
	moves    []map[string]interface{}
}

func NewTicTacToeChecker() *TicTacToeChecker {
	ret := &TicTacToeChecker{
		moves: make([]map[string]interface{}, 0),
	}
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			ret.gameGrid[i][j] = EMPTY_CELL
		}
	}
	return ret
}

//IsValidMove checks if a move is valid.
func (receiver *TicTacToeChecker) IsValidMove(Move map[string]interface{}) bool {
	Cell := Move["Cell"].(map[string]interface{})
	x := Cell["X"].(int)
	y := Cell["Y"].(int)
	return receiver.gameGrid[x][y] == EMPTY_CELL
}

//MakeMove does a move, if it is valid.
func (receiver *TicTacToeChecker) MakeMove(Move map[string]interface{}) (bool, Result) {
	valid := receiver.IsValidMove(Move)
	if valid {
		Cell := Move["Cell"].(map[string]interface{})
		x := Cell["X"].(int)
		y := Cell["Y"].(int)
		symbol := Move["Symbol"].(TicTacToeSymbol)
		receiver.gameGrid[x][y] = symbol
	}
	receiver.moves = append(receiver.moves, Move)
	won := receiver.CheckWin(Move)
	return valid, won
}

func (receiver *TicTacToeChecker) CheckWin(lastMove map[string]interface{}) Result {
	Cell := lastMove["Cell"].(map[string]interface{})
	x := Cell["X"].(int)
	y := Cell["Y"].(int)
	symbol := lastMove["Symbol"].(TicTacToeSymbol)
	//check cross
	//check rows
	for i := 0; i < 3; i++ {
		if receiver.gameGrid[x][i] != symbol {
			break
		}
		if i == 2 {
			return WIN
		}
	}
	//check columns
	for i := 0; i < 3; i++ {
		if receiver.gameGrid[i][y] != symbol {
			break
		}
		if i == 2 {
			return WIN
		}
	}
	//check diags
	//diag
	if x == y {
		for i := 0; i < 3; i++ {
			if receiver.gameGrid[i][i] != symbol {
				break
			}
			if i == 2 {
				return WIN
			}
		}
	}
	//anti diag
	if x+y == 2 {
		for i := 0; i < 3; i++ {
			if receiver.gameGrid[i][2-i] != symbol {
				break
			}
			if i == 2 {
				return WIN
			}
		}
	}
	if len(receiver.moves) == 9 {
		return DRAW
	}
	return ONGOING
}
