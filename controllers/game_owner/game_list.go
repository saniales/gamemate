package gameOwnerController

import (
	"sanino/gamemate/configurations"
	"sanino/gamemate/models/game_owner/data_structures"
)

//GetGames gets the enabled games for the owner (from the archives)
//Will be cached by the client (owner app).
//Version 1.0
//QUESTION: does it have sense to add cache even here?
func GetGames(ownerID int64) ([]gameOwnerDataStructs.Game, error) {
	/*
		  games, err := getGamesFromCache(ownerID)
			if err != nil {
				games, err = getGamesFromArchives(ownerID)
				if err != nil {
					return nil, false, err
				}
				err = updateCacheWithGameList(ownerID, games)
				if err != nil {
					return games, false, err
				}
				return games, true, err
			}
			return games, true, nil
	*/
	return getGamesFromArchives(ownerID)
}

//getGamesFromArchives gets from the archives the games of a particular owner.
func getGamesFromArchives(ownerID int64) ([]gameOwnerDataStructs.Game, error) {
	stmtQuery, err := configurations.ArchivesPool.Prepare("SELECT gameID, name, description, maxPlayers FROM games WHERE ownerID = ?")
	if err != nil {
		return nil, err
	}
	defer stmtQuery.Close()

	rows, err := stmtQuery.Query(ownerID)
	if err != nil {
		return nil, err
	}

	ret := make([]gameOwnerDataStructs.Game, 0, 10)

	for !rows.Next() {
		game := gameOwnerDataStructs.Game{}
		err = rows.Err()
		if err != nil {
			return nil, err
		}
		err = rows.Scan(&game.ID, &game.Name, &game.Description, &game.MaxPlayers)
		if err != nil {
			return nil, err
		}
		//get game stats?
		ret = append(ret, game)
	}

	return ret, nil
}
