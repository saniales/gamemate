package gameOwnerController

import (
	"errors"
	"fmt"

	"sanino/gamemate/configurations"
	"sanino/gamemate/models/game_owner/data_structures"
)

//updateCacheWithNewGame updates the Cache with the specified API_Token.
//
//Return error if did not manage to update the cache.
func updateCacheWithNewGame(Game gameOwnerDataStructs.Game) error {
	conn := configurations.CachePool.Get()
	defer conn.Close()

	err := conn.Send("MULTI")
	if err != nil {
		return err
	}

	//a game must be in cache until it's removed
	err = conn.Send("SADD", "all_games", Game.ID)
	if err != nil {
		return err
	}

	err = conn.Send("HMSET", fmt.Sprintf("games/with_id/%d", Game.ID), "name", Game.Name, "description", Game.Description, "max_players", Game.MaxPlayers)
	if err != nil {
		return err
	}

	err = conn.Send("EXEC")
	if err != nil {
		return err
	}
	err = conn.Flush()
	if err != nil {
		return err
	}
	return nil
}

//removeGameFromCache removes from the Cache the specified API_Token.
//
//Return error if did not manage to update the cache.
func removeGameFromCache(gameID int64) error {
	conn := configurations.CachePool.Get()
	defer conn.Close()

	err := conn.Send("MULTI")
	if err != nil {
		return err
	}

	err = conn.Send("SREM", "games", gameID)
	if err != nil {
		return err
	}

	err = conn.Send("DEL", fmt.Sprintf("games/with_id/%d", gameID))
	if err != nil {
		return err
	}

	err = conn.Flush()
	if err != nil {
		return err
	}

	err = conn.Send("EXEC")
	if err != nil {
		return err
	}

	return nil
}

//checkGameInArchives checks for the existance of the game in the archives
//and, if found, updates the cache.
//
//Return true if found, false otherwise.
func checkGameInArchives(name string, ownerID int64) (bool, error) {
	//check if present
	stmtQuery, err := configurations.ArchivesPool.Prepare("SELECT COUNT(ID), ID FROM games WHERE name = ? AND ownerID = ? GROUP BY ID")
	if err != nil {
		return false, err
	}
	defer stmtQuery.Close()

	result, err := stmtQuery.Query(name, ownerID)
	if err != nil {
		return false, err
	}
	if !result.Next() {
		return false, errors.New("Check game error (archives) : Empty Table, Query with errors (should report 0 when item is not in table)")
	}

	var num_rows int64
	var gameID int64
	result.Scan(&num_rows, gameID)

	if num_rows > 0 {
		//gets full game from ID
		stmtQuery, err = configurations.ArchivesPool.Prepare(fmt.Sprintf("SELECT name, description, max_players FROM games WHERE ID = %d", gameID))
		if err != nil {
			return false, err
		}
		defer stmtQuery.Close()

		result, err = stmtQuery.Query()
		if err != nil {
			return false, err
		}
		if !result.Next() {
			return false, errors.New("Check game error (archives) : Empty Table, Query with errors")
		}
		var game gameOwnerDataStructs.Game
		game.ID = gameID
		result.Scan(&game.Name, &game.Description, &game.MaxPlayers)
		err = updateCacheWithNewGame(game)
		if err != nil { //did not update cache but the request has been satisfied.
			return true, nil
		}
		return true, nil
	}
	return false, nil
}

//addGameInArchives adds a game linked to the specified owner to the archives.
//
//return insertID of the game if successfull, otherwise fills error.
func addGameInArchives(ownerID int64, name string, description string, maxPlayers int64) (int64, error) {
	//TODO: find a way to handle duplicates. or leave the query fail and retry.
	stmtQuery, err := configurations.ArchivesPool.Prepare("INSERT INTO games (ID, owner_ID, name, description, max_players) VALUES (NULL, ?, ?, ?, ?)")
	if err != nil {
		return -1, err
	}
	defer stmtQuery.Close()
	result, err := stmtQuery.Exec(ownerID, name, description, maxPlayers)
	if err != nil {
		return -1, err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return -1, err
	}
	if rows <= 0 {
		return -1, errors.New("No Row Affected, possible problem with the query")
	}

	insertID, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}
	return insertID, nil
}

//removeGameFromArchives removes a token from the Archives.
func removeGameFromArchives(ownerID int64, gameID int64) error {
	stmtQuery, err := configurations.ArchivesPool.Prepare("DELETE FROM games WHERE ID = ? AND owner_ID = ?")
	if err != nil {
		return err
	}
	defer stmtQuery.Close()

	result, err := stmtQuery.Exec(gameID, ownerID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows <= 0 {
		return errors.New("No Row Affected, possible problem with the query, or the owner is fake")
	}
	return nil
}
