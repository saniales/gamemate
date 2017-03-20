package outGameController

import (
	"errors"
	"fmt"
	"sanino/gamemate/configurations"
	"sanino/gamemate/constants"
	"sanino/gamemate/models/user/data_structures"
	"strconv"

	"github.com/garyburd/redigo/redis"
)

//getGamesFromCache gets the games from the cache, if present and consistent.
func getGamesFromCache(userID int64) ([]userDataStructs.SummarizedGame, error) {
	return nil, errors.New("Currently not supported")
}

//getGamesFromArchives gets the games from the archives, WITHOUT Currently playing
//users (which are in cache).
//NOTE: if update the cache with this values remember to update call updateCacheCurrentlyPlayed().
func getGamesFromArchives(userID int64) ([]userDataStructs.SummarizedGame, error) {
	Query := fmt.Sprintf(
		"SELECT g.gameID, g.name, 1 "+
			"FROM games g JOIN user_game_enabled u "+
			"ON u.gameID = g.gameID AND u.userID = %d "+
			"UNION "+
			"SELECT g2.gameID, g2.name, 0 "+
			"FROM games g2 LEFT JOIN user_game_enabled u2 "+
			"ON u2.gameID = g2.gameID AND u2.userID IS NULL",
		userID,
	)
	stmtQuery, err := configurations.ArchivesPool.Prepare(Query)
	if err != nil {
		return nil, err
	}

	result, err := stmtQuery.Query()
	if err != nil {
		return nil, err
	}

	ret := make([]userDataStructs.SummarizedGame, 0, 10)
	for !result.Next() {
		var ID int64
		var name string
		var enabledInt int64
		result.Scan(&ID, &name, &enabledInt)
		game := userDataStructs.SummarizedGame{ID: ID, Name: name, Enabled: enabledInt > 0}
		ret = append(ret, game)
	}

	err = result.Err()
	if err != nil {
		return nil, err
	}

	return ret, nil
}

//GetGames get the games from a requesting user, setting enabled flag.
//NOTE: incomplete, do not use cache without finishing function implementation.
func GetGames(userID int64) ([]userDataStructs.SummarizedGame, bool, error) {
	//games, err := getGamesFromCache(userID)
	//if err != nil {
	games, err := getGamesFromArchives(userID)
	if err != nil {
		return nil, false, err
	}
	return games, false, nil
	//err = updateCacheAllGames(userID)
	//if err != nil {
	//return games, false, nil
	//}
	//return games, true, nil
	//}
	//return games, false, nil
}

func getEnabledGameIDsFromArchives(userID int64) ([]int64, error) {
	stmtQuery, err := configurations.ArchivesPool.Prepare("SELECT gameID FROM user_game_enabled WHERE userID = ?")
	if err != nil {
		return nil, err
	}
	defer stmtQuery.Close()

	result, err := stmtQuery.Query(userID)
	if err != nil {
		return nil, err
	}
	defer result.Close()

	ret := make([]int64, 0, 10)

	for !result.Next() {
		var gameID int64

		err = result.Err()
		if err != nil {
			return nil, err
		}

		err = result.Scan(&gameID)
		if err != nil {
			return nil, err
		}
		ret = append(ret, gameID)
	}
	return ret, nil
}

//GetEnabledGameIDs gets the IDs of the games enabled for a user.
func GetEnabledGameIDs(userID int64) ([]int64, bool, error) {
	games, err := getEnabledGameIDsFromCache(userID)
	if err != nil {
		games, err = getEnabledGameIDsFromArchives(userID)
		if err != nil {
			return nil, false, err
		}
		return games, false, err
	}
	return games, true, nil
}

func getEnabledGameIDsFromCache(userID int64) ([]int64, error) {
	conn := configurations.CachePool.Get()
	defer conn.Close()

	strings, err := redis.Strings(conn.Do("SMEMBERS", fmt.Sprintf(constants.USER_ENABLED_GAMES, userID)))
	if err != nil {
		return nil, err
	}

	ints := make([]int64, len(strings))
	for i, v := range strings {
		ints[i], err = strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, err
		}
	}
	return ints, nil
}

//UpdateCacheAllGames updates the cache with all games summarized data.
//
//NOTE: Complexity O(N log(N)), so should be used as least as possible.
//
//N => number of games to add.
func UpdateCacheAllGames(games []userDataStructs.SummarizedGame) error {
	return errors.New("Currently not supported")
}

func getGameDetail(gameID int64, userID int64) (userDataStructs.SummarizedGame, bool, error) {
	game, err := getGameDetailFromCache(gameID, userID)
	if err != nil {
		game, err = getGameDetailFromArchives(gameID, userID)
		if err != nil {
			return userDataStructs.SummarizedGame{}, false, err
		}
		err := updateCacheGameDetail(game, userID)
		if err != nil {
			return game, false, nil
		}
		return game, true, nil
	}
	return game, false, nil
}

func getGameDetailFromCache(gameID int64, userID int64) (userDataStructs.SummarizedGame, error) {
	return userDataStructs.SummarizedGame{}, errors.New("Currently not supported")
}

func getGameDetailFromArchives(gameID int64, userID int64) (userDataStructs.SummarizedGame, error) {
	ret := userDataStructs.SummarizedGame{}
	QueryGetDetail := fmt.Sprintf(
		"SELECT gameID, name, maxPlayers FROM games WHERE gameID = %d",
		gameID,
	)

	stmtQuery, err := configurations.ArchivesPool.Prepare(QueryGetDetail)
	if err != nil {
		return ret, err
	}

	var ID int64
	var name string
	err = stmtQuery.QueryRow(0).Scan(&ID, &name)
	if err != nil {
		return ret, err
	}

	stmtQuery.Close()

	QueryGetEnabled := fmt.Sprintf(
		"SELECT COUNT(*) FROM user_game_enabled WHERE gameID = %d AND userID = %d",
		gameID, userID,
	)

	stmtQuery, err = configurations.ArchivesPool.Prepare(QueryGetEnabled)
	if err != nil {
		return ret, err
	}
	defer stmtQuery.Close()

	var count int64
	err = stmtQuery.QueryRow(0).Scan(&count)
	if err != nil {
		return ret, err
	}

	ret.ID = gameID
	ret.Name = name
	ret.Enabled = count > 0

	return ret, nil
}

func updateCacheGameDetail(game userDataStructs.SummarizedGame, userID int64) error {
	return errors.New("Currently not supported")
}
