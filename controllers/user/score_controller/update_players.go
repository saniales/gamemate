package scoreController

import (
	"fmt"
	"sanino/gamemate/configurations"
	"sanino/gamemate/constants"

	"github.com/garyburd/redigo/redis"
)

//GetWeeklyPlayers gets the number of players who have joined a match in the last week.
func GetWeeklyPlayers(gameID int64) (int64, error) {
	return getPlayingPlayersInPeriod(gameID, "week")
}

//GetMonthlyPlayers gets the number of players who have joined a match in the last month.
func GetMonthlyPlayers(gameID int64) (int64, error) {
	return getPlayingPlayersInPeriod(gameID, "month")
}

//GetYearlyPlayers gets the number of players who have joined a match in the last year.
func GetYearlyPlayers(gameID int64) (int64, error) {
	return getPlayingPlayersInPeriod(gameID, "year")
}

//getPlayingPlayersInPeriod gets the number of players of a game in the period of time (cache only)
func getPlayingPlayersInPeriod(gameID int64, period string) (int64, error) {
	conn := configurations.CachePool.Get()
	defer conn.Close()

	return redis.Int64(conn.Do("ZCARD", fmt.Sprintf("%s/%d:%s_players", constants.GAMES_SET, gameID, period)))
}
