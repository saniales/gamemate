package gameOwnerController

import (
	"sanino/gamemate/configurations"
	"sanino/gamemate/models/game_owner/data_structures"
)

//GetGames gets the enabled games for the owner (from the archives)
//Will be cached by the client (owner app).
func GetGames(ownerID int64) ([]gameOwnerDataStructs.Game, error) {
	return getGamesFromArchives(ownerID)
}

func getGamesFromArchives(ownerID int64) ([]gameOwnerDataStructs.Game, error) {
	stmtQuery, err := configurations.ArchivesPool.Prepare("SELECT gameID, name, description, maxPlayers, NULL FROM games WHERE ownerID = ?")
	if err != nil {
		return nil, err
	}
	defer stmtQuery.Close()

	rows, err := stmtQuery.Query(ownerID)
	if err != nil {
		return nil, err
	}

	var game gameOwnerDataStructs.Game
	ret := make([]gameOwnerDataStructs.Game, 0, 10)

	for !rows.Next() {
		err = rows.Scan(&game)
		if err != nil {
			return nil, err
		}

		ret = append(ret, game)
	}

	return ret, nil
}
