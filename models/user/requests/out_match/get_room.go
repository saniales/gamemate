package outMatchRequests

import (
	"errors"
)

type GetRoom struct {
	//UserID int64
	//GameID   int64  `json:"GameID" xml:"GameID" form:"GameID"`
	Username string `json:"Username" xml:"Username" form:"Username"` //prototypical, to avoid use cache for faster development //TODO: add cache layer
}

//FromMap parses the a map and creates a valid struct, or
//returns an error.
func (receiver *GetRoom) FromMap(Map map[string]interface{}) error {
	//userID, err := strconv.ParseInt(Map["UserID"], 10, 64)
	//if err != nil {
	//return errors.New("Invalid Map : UserID not found")
	//}
	//gameID, err := strconv.ParseInt(Map["GameID"], 10, 64)
	//if err != nil || gameID <= 0 {
	//return errors.New("Invalid Map : GameID not found")
	//}
	tmpUsername := Map["Username"].(string)
	if tmpUsername == "" {
		return errors.New("Invalid Map : Username not found")
	}
	//receiver.UserID = userID
	//receiver.gameID = gameID
	receiver.Username = tmpUsername
	return nil
}
