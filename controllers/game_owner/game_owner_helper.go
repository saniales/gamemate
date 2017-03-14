package gameOwnerController

import (
	"errors"
	"fmt"

	"github.com/garyburd/redigo/redis"

	"sanino/gamemate/configurations"
	"sanino/gamemate/constants"
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
	err = conn.Send("ZADD", constants.SUMMARY_GAMES_SET, 0, fmt.Sprintf("%d:%s", Game.ID, Game.Name))
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

	_, err = conn.Do("EXEC")
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
	defer result.Close()

	if !result.Next() {
		return false, errors.New("Check game error (archives) : Empty Table, Query with errors (should report 0 when item is not in table)")
	}

	var num_rows int64
	var gameID int64
	result.Scan(&num_rows, &gameID)

	if num_rows > 0 {
		//gets full game from ID
		stmtQuery, err = configurations.ArchivesPool.Prepare(fmt.Sprintf("SELECT name, description, maxPlayers FROM games WHERE gameID = %d", gameID))
		if err != nil {
			return false, err
		}
		defer stmtQuery.Close()

		result, err = stmtQuery.Query()
		if err != nil {
			return false, err
		}
		defer result.Close()

		if !result.Next() { //not found
			return false, nil
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
	stmtQuery, err := configurations.ArchivesPool.Prepare("INSERT INTO games (gameID, ownerID, name, description, maxPlayers) VALUES (NULL, ?, ?, ?, ?)")
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

//EnableDisableGameForUser makes an action on the specified game for the user, if valid.
//
//returns if the cache has been updated and the error if present.
func EnableDisableGameForUser(userID int64, gameID int64, enable bool) (bool, error) {
	err := enableDisableGameInArchives(userID, gameID, enable)
	if err != nil {
		return false, err
	}
	//updates cache
	if enable {
		err = enableGameInCache(userID, gameID)
		if err != nil {
			return false, nil
		}
		return true, nil
	}
	return true, err
}

//enableGameInCache enables a game for the user; it acts on the cache layer.
func enableGameInCache(userID int64, gameID int64) error {
	conn := configurations.CachePool.Get()
	defer conn.Close()

	return conn.Send("SADD", fmt.Sprintf("games/with_id/%d:enabled_players", gameID), userID)
}

//EnableDisableGameInArchives enables (or disables) a game for a user.
func enableDisableGameInArchives(userID int64, gameID int64, enable bool) error {
	var query string
	if enable {
		query = "INSERT INTO user_game_enabled (userID, gameID) VALUES (?, ?)"
	} else {
		query = "DELETE FROM user_game_enabled WHERE userID = ? and gameID = ?"
	}

	stmtQuery, err := configurations.ArchivesPool.Prepare(query)
	if err != nil {
		return err
	}
	defer stmtQuery.Close()

	result, err := stmtQuery.Exec(userID, gameID)
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

//GetOwnerOfGame return the owner ID of the specified game.
func GetOwnerOfGame(gameID int64) (int64, bool, error) {
	ownerID, err := getOwnerFromGameInCache(gameID)
	if err != nil {
		ownerID, err = getOwnerFromArchives(gameID)
		if err != nil {
			return -1, false, err
		}
		err = updateCacheOwnerOfGame(gameID, ownerID)
		if err != nil {
			return ownerID, false, nil
		}
		return ownerID, true, nil
	}
	return ownerID, true, nil
}

//getOwnerFromGameInCache get the owner id from a game, looking for it in the cache.
func getOwnerFromGameInCache(gameID int64) (int64, error) {
	conn := configurations.CachePool.Get()
	defer conn.Close()

	return redis.Int64(conn.Do("HMGET", fmt.Sprintf("games/with_id/%d", gameID), "owner_id"))
}

//updateCacheOwnerOfGame updates the cache, setting owner for the specified game.
func updateCacheOwnerOfGame(gameID int64, ownerID int64) error {
	conn := configurations.CachePool.Get()
	defer conn.Close()

	err := conn.Send("HMSET", fmt.Sprintf("games/with_id/%d", gameID), "owner_id", ownerID)
	if err != nil {
		return err
	}
	err = conn.Flush()
	if err != nil {
		return err
	}
	return nil
}

//getOwnerFromArchives returns the  owner of the specified game, if exists; it
//searches in the archives.
func getOwnerFromArchives(gameID int64) (int64, error) {
	stmtQuery, err := configurations.ArchivesPool.Prepare("SELECT ownerID FROM games WHERE gameID = ?")
	if err != nil {
		return -1, err
	}
	defer stmtQuery.Close()

	result, err := stmtQuery.Query(gameID)
	if err != nil {
		return -1, err
	}
	defer result.Close()

	if !result.Next() {
		return -1, nil //not found
	}

	err = result.Err()
	if err != nil {
		return -1, err
	}
	var ownerID int64
	err = result.Scan(&ownerID)
	if err != nil {
		return -1, err
	}
	return ownerID, nil
}
