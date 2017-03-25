package userDataStructs

//Room represents a lobby where players wait to enter a match.
type Room struct {
	ID           int64    `json:"ID" xml:"ID" form:"ID"`
	Players      []Player `json:"Players" xml:"Players" form:"Players"`
	MatchStarted bool     `json:"MatchStarted" xml:"MatchStarted" form:"MatchStarted"`
}
