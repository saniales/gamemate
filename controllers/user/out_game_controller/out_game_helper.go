package outGameController

import (
	"fmt"
	"sanino/gamemate/configurations"
	"sanino/gamemate/constants"
	"sanino/gamemate/models/user/data_structures"
	"strconv"

	"github.com/garyburd/redigo/redis"
)

//getGamesFromCache gets the games from the cache, if present and consistent.
func getGamesFromCache() ([]userDataStructs.SummarizedGame, error) {
	conn := configurations.CachePool.Get()
	defer conn.Close()

	values, err := redis.Strings(conn.Do("ZRANGE", constants.SUMMARY_GAMES_SET, "0", "-1", "WITHSCORES"))
	if err != nil {
		return nil, err
	}

	ret := make([]userDataStructs.SummarizedGame, 0, 10)
	var tempGame userDataStructs.SummarizedGame

	for i, value := range values {
		if i%2 == 0 { //object
			tempGame = userDataStructs.SummarizedGame{Name: value, CurrentlyPlayedBy: -1}
		} else { //score
			tempGame.CurrentlyPlayedBy, err = strconv.ParseInt(value, 10, 64)
			if err != nil {
				return nil, err
			}
			ret = append(ret, tempGame)
		}
	}
	//tries to
	conn.Do("EXPIRE", constants.SUMMARY_GAMES_SET, constants.CACHE_REFRESH_INTERVAL)
	return ret, nil
}

//getGamesFromArchives gets the games from the archives, WITHOUT Currently playing
//users (which are in cache).
//NOTE: if update the cache with this values remember to update call updateCacheCurrentlyPlayed().
func getGamesFromArchives() ([]userDataStructs.SummarizedGame, error) {
	stmtQuery, err := configurations.ArchivesPool.Prepare("SELECT gameID, name FROM games")
	if err != nil {
		return nil, err
	}

	result, err := stmtQuery.Query()
	if err != nil {
		return nil, err
	}

	ret := make([]userDataStructs.SummarizedGame, 0, 10)
	for !result.Next() {
		game := userDataStructs.SummarizedGame{CurrentlyPlayedBy: 0}
		result.Scan(&game.ID, &game.Name)
		ret = append(ret, game)
	}
	if result.Err() != nil {
		return nil, result.Err()
	}

	return ret, nil
}

func getEnabledGameIDsFromArchives(userID int64) ([]int64, error) {
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

//GetGames get summarized data for games.
func GetGames() ([]userDataStructs.SummarizedGame, bool, error) {
	games, err := getGamesFromCache()
	if err != nil {
		games, err = getGamesFromArchives()
		if err != nil {
			return nil, false, err
		}
		err = UpdateCacheAllGames(games)
		return games, err != nil, err
	}
	return games, true, nil
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
	conn := configurations.CachePool.Get()
	defer conn.Close()

	err := conn.Send("MULTI")
	if err != nil {
		return err
	}

	for _, game := range games {
		err = conn.Send("ZADD", constants.SUMMARY_GAMES_SET, 0, fmt.Sprintf("%d:%s", game.ID, game.Name))
		if err != nil {
			return err
		}
	}

	_, err = conn.Do("EXEC")
	return err
}
