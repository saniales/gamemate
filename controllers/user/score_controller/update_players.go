package scoreController

import (
	"fmt"
	"sanino/gamemate/configurations"

	"github.com/garyburd/redigo/redis"
)

//GetWeeklyPlayers gets the number of players who have joined a match in the last week.
func GetWeeklyPlayers(gameID int64) (int64, error) {
	return getPlayingPlayersInPeriod(gameID, "weekly")
}

//GetMonthlyPlayers gets the number of players who have joined a match in the last month.
func GetMonthlyPlayers(gameID int64) (int64, error) {
	return getPlayingPlayersInPeriod(gameID, "monthly")
}

//GetYearlyPlayers gets the number of players who have joined a match in the last year.
func GetYearlyPlayers(gameID int64) (int64, error) {
	return getPlayingPlayersInPeriod(gameID, "yearly")
}

//getPlayingPlayersInPeriod gets the number of players of a game in the period of time (cache only)
func getPlayingPlayersInPeriod(gameID int64, period string) (int64, error) {
	conn := configurations.CachePool.Get()
	defer conn.Close()

	return redis.Int64(conn.Do("ZCARD", fmt.Sprintf("games/with_id/%d:%s_players", gameID, period)))
}
