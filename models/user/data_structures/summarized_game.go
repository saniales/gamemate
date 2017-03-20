package userDataStructs

import "strconv"

//SummarizedGame represents a game for the front end without details, only operational
//and essentiald data.
type SummarizedGame struct {
	ID   int64  `json:"ID" xml:"ID" form:"ID"`       //The ID of the game.
	Name string `json:"Name" xml:"Name" form:"Name"` //The Name of the game
	//CurrentlyPlayedBy int64  `json:"CurrentlyPlayedBy" xml:"CurrentlyPlayedBy" form:"CurrentlyPlayedBy"` //The number of logged players who plays currnetly this game.
	Enabled bool `json:"Enabled" xml:"Enabled" form:"Enabled"` //True if the game is enabled for the user, false otherwise.
}

//FromMap gets the values to fill the struct from a map passed as parameter.
func (receiver *SummarizedGame) FromMap(Map map[string]string) error {
	ID, err := strconv.ParseInt(Map["ID"], 10, 64)
	if err != nil {
		return err
	}
	receiver.ID = ID
	receiver.Name = Map["Name"]
	receiver.Enabled, err = strconv.ParseBool(Map["Enabled"])
	if err != nil {
		return err
	}
	return nil
}
