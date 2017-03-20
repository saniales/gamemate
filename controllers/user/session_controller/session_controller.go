package sessionController

import (
	"sanino/gamemate/constants"
	"sanino/gamemate/controllers/shared"
)

//GetUserIDFromSessionToken gets the user ID from the session token in cache.
func GetUserIDFromSessionToken(token string) (int64, error) {
	return controllerSharedFuncs.GetIDFromSessionSet(constants.LOGGED_USERS_SET, token)
}
