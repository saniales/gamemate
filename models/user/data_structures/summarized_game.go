package userDataStructs

//SummarizedGame represents a game for the front end without details, only operational
//and essentiald data.
type SummarizedGame struct {
	ID                int64  `json:"ID" xml:"ID" form:"ID"`                                              //The ID of the game.
	Name              string `json:"Name" xml:"Name" form:"Name"`                                        //The Name of the game
	CurrentlyPlayedBy int64  `json:"CurrentlyPlayedBy" xml:"CurrentlyPlayedBy" form:"CurrentlyPlayedBy"` //The number of logged players who plays currnetly this game.
}
