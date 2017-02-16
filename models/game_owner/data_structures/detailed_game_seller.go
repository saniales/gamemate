package gameOwnerDataStructs

//Game represents a single game from the seller Point of View Additional Data.
type Game struct {
	ID               int64  `json:"ID" xml:"ID" form:"ID"`                                           //ID of the Game.
	Name             string `json:"Name" xml:"Name" form:"Name"`                                     //Name of the Game.
	Description      string `json:"Description" xml:"Description" form:"Description"`                //Description of the game.
	MaxPlayers       int64  `json:"MaxPlayers" xml:"MaxPlayers" form:"MaxPlayers"`                   //Max Players per match for this game.
	ActiveUsersYear  int64  `json:"ActiveUsersYear" xml:"ActiveUsersYear" form:"ActiveUsersYear"`    //number of active users in this year
	ActiveUsersMonth int64  `json:"ActiveUsersMonth" xml:"ActiveUsersMonth" form:"ActiveUsersMonth"` //number of active users in this month
	ActiveUsersWeek  int64  `json:"ActiveUsersWeek" xml:"ActiveUsersWeek" form:"ActiveUsersWeek"`    //number of active users in this week
	ActiveUsersToday int64  `json:"ActiveUsersToday" xml:"ActiveUsersToday" form:"ActiveUsersToday"` //number of active users today
	AvgRatingYear    int16  `json:"AvgRatingYear" xml:"AvgRatingYear" form:"AvgRatingYear"`          //avg rating for users in one year.
}
