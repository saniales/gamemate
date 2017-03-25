package userDataStructs

import (
	"errors"
	"strconv"
)

//Player represents a playing user.
type Player struct {
	ID       int64  `json:"ID" xml:"ID" form:"ID"`                   //The user ID of the player.
	Username string `json:"Username" xml:"Username" form:"Username"` //The username of the player.
}

func (receiver *Player) FromMap(Map map[string]string) error {
	ID, err := strconv.ParseInt(Map["ID"], 10, 64)
	if err != nil {
		return err
	}
	receiver.ID = ID
	tmpUsername := Map["Username"]
	if tmpUsername == "" {
		return errors.New("Player:FromMap error : Username field not found in Map")
	}
	receiver.Username = tmpUsername
	return nil
}
