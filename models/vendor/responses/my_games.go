package vendorResponses

//MyGames represents a response from the server to a vendorResponses.MyGames
type MyGames struct {
	Type string `json:"Type" xml:"Type" form:"Type"`
	//Games []Game `json:"Games" xml:"Games" form:"Games"`
}

//TODO : func FromGames([]Game) { /* implement */ }
