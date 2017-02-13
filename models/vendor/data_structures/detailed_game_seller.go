package vendorDataStructs

//Game represents a single game from the seller Point of View.
type Game struct {
	ID          int64
	Name        string
	Description string
	MaxPlayers  int64
}
