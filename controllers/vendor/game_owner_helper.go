package gameOwnerController

import (
	"errors"

	"sanino/gamemate/configurations"
	"sanino/gamemate/controllers/shared"
	"sanino/gamemate/models/game_owner/data_structures"
)

//updateCacheWithNewGame updates the Cache with the specified API_Token.
//
//Return error if did not manage to update the cache.
func updateCacheWithNewGame(game gameOwnerDataStructs.Game) error {
	conn := configurations.CachePool.Get()
	defer conn.Close()

	err := conn.Send("MULTI")
	if err != nil {
		return err
	}

	//a game must be in cache until it's removed
	err = conn.Send("SADD", "all_games", Game.Name)
	if err != nil {
		return err
	}

	err = conn.Send("HMSET", "games/"+Game.Name, "ID", Game.ID, "max_players", Game.MaxPlayers)
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
func removeGameFromCache(Name string) error {
	conn := configurations.CachePool.Get()
	defer conn.Close()

	err := conn.Send("MULTI")
	if err != nil {
		return err
	}

	err = conn.Send("SREM", "all_games", Name)
	if err != nil {
		return err
	}

	err = conn.Send("DEL", "games/"+Name)
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
	stmtQuery, err := configurations.ArchivesPool.Prepare("SELECT COUNT(name) FROM games WHERE name = ? AND ownerID = ?")
	if err != nil {
		return false, err
	}
	defer stmtQuery.Close()
	result, err := stmtQuery.Query(name, ownerID)
	if err != nil {
		return false, err
	}
	if !result.Next() {
		return false, errors.New("Check API error (archives) : Empty Table, Query with errors (should report 0 when item is not in table)")
	}

	var num_rows int64
	result.Scan(&num_rows)

	if num_rows > 0 {
		err = updateCacheWithNewGame(name)
		if err != nil { //did not update cache but the request has been satisfied.
			return true, nil
		}
		return true, nil
	}
	return false, nil
}

//addAPI_TokenInArchives adds a token linked to the specified developer to the archives.
func addGameInArchives(ownerID int64) (string, error) {
	token := controllerSharedFuncs.GenerateToken()
	//TODO: find a way to handle duplicates. or leave the query fail and retry.
	stmtQuery, err := configurations.ArchivesPool.Prepare("INSERT INTO games (ownerId, name, description, max_players) VALUES (?, ?, ?)")
	if err != nil {
		return "", err
	}
	defer stmtQuery.Close()
	result, err := stmtQuery.Exec(developerEmail, token)
	if err != nil {
		return "", err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return "", err
	}
	if rows <= 0 {
		return "", errors.New("No Row Affected, possible problem with the query")
	}
	return token, nil
}

//removeGameFromArchives removes a token from the Archives.
func removeGameFromArchives( /*Owner string, */ Name string) error {
	stmtQuery, err := configurations.ArchivesPool.Prepare("DELETE FROM games WHERE Name = ?")
	if err != nil {
		return err
	}
	defer stmtQuery.Close()

	result, err := stmtQuery.Exec(Name)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows <= 0 {
		return errors.New("No Row Affected, possible problem with the query")
	}
	return nil
}
