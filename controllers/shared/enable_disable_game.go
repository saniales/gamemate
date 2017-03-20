package controllerSharedFuncs

import (
	"fmt"
	"sanino/gamemate/configurations"
)

//EnableDisableGameForUser makes an action on the specified game for the user, if valid.
//
//returns if the cache has been updated and the error if present.
func EnableDisableGameForUser(userID int64, gameID int64, enable bool) (bool, error) {
	err := enableDisableGameInArchives(userID, gameID, enable)
	if err != nil {
		return false, err
	}
	err = enableDisableGameInCache(userID, gameID, enable)
	if err != nil {
		return false, nil
	}
	return true, nil
}

//enableDisableGameInCache enables or disables a game for the user; it acts on the cache layer.
func enableDisableGameInCache(userID int64, gameID int64, enable bool) error {
	conn := configurations.CachePool.Get()
	defer conn.Close()
	command := "SREM"
	if enable {
		command = "SADD"
	}
	return conn.Send(command, fmt.Sprintf("games/with_id/%d:enabled_players", gameID), userID)
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
	if rows <= 0 { //no row affected
		return nil
	}
	return nil
}
