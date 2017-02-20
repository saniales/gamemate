package gameOwnerController

import "sanino/gamemate/configurations"

//GetGamesEnabled gets the enabled games for the user.
func GetGamesEnabled(userID int64) ([]Game, bool, error) {
	Games, err := getGamesEnabledFromCache(userID)
	if err != nil {
		Games, err = getGamesEnabledFromArchives(userID)
		if err != nil {
			return nil, false, err
		}
		return Games, false, nil
	}
	return Games, true, nil
}

func getGamesEnabledFromCache(gameID int64) ([]Game, error) {
	conn := configurations.CachePool.Get()
	defer conn.Close()
}
