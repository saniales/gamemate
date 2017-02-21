package outGameController

import (
	"sanino/gamemate/configurations"
	"sanino/gamemate/models/user/data_structures"
)

func getGamesFromCache() error {
	conn := configurations.CachePool.Get()
	defer conn.Close()

	panic("IMPLEMENT")
}

func getGamesFromArchives(userID int64) {

}

func getEnabledGamesFromCache(userID int64) {

}

func getEnabledGamesFromArchives(userID int64) {

}

//GetGames get summarized data for games.
func GetGames() ([]userDataStructs.SummarizedGame, error) {
	panic("IMPLEMENT")
}

//GetEnabledGameIDs gets the IDs of the games enabled for a user.
func GetEnabledGameIDs(userID) ([]int64, error) {
	panic("IMPLEMENT")
}
