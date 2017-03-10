package constants

import "time"

const (
	ROOT_PATH                   string = "/"                  //ROOT_PATH is the path to the root directory of the server.
	USER_REGISTRATION_PATH      string = "/user/register"     //USER_REGISTRATION_PATH is the path to handle user registration.
	AUTH_PATH                   string = "/user/auth"         //AUTH_PATH is the path to handle user authentication.
	USER_ALL_GAME_LIST_PATH     string = "/games/get/all"     //USER_ALL_GAME_LIST_PATH is the path to handle a request to the list of games
	USER_ENABLED_GAME_LIST_PATH string = "/games/get/enabled" //USER_ENABLED_GAME_LIST_PATH is the path to handle a request to get the games enabled for a user.
	ROOM_CREATION_PATH          string = "/room/create"       //ROOM_CREATION_PATH is the path to handle the creation of a new room.
	GET_ROOM_REQUEST_PATH       string = "/room/get"          //GET_ROOM_REQUEST_PATH is the path to handle the request of data of a particular room.

	MATCH_CREATION_PATH     string = "/match/create" //MATCH_CREATION_PATH is the path to handle the creationd of a match (not started).
	MATCH_START_PATH        string = "/match/start"  //MATCH_START_PATH is the path to handle the start of a match (it becomes LIVE).
	MATCH_DATA_REQUEST_PATH string = "/match/get"    //MATCH_DATA_REQUEST_PATH is the path to handle the request of getting data of a particular match.

	TURN_ACTION_PATH string = "/user/move/make"     //TURN_ACTION_PATH is the path to handle an action in a match.
	TURN_END_PATH    string = "/user/move/end_turn" //TURN_END_PATH is the path to handle the end of a turn.

	DEVELOPER_AUTH_PATH           string = "/dev/auth"           //DEVELOPER_AUTH_PATH is the path to handle developer authentication.
	DEVELOPER_REGISTRATION_PATH   string = "/dev/register"       //DEVELOPER_REGISTRATION_PATH is the path to hande developer registration.
	DEVELOPER_ADD_API_TOKEN_PATH  string = "/dev/api_token/add"  //DEVELOPER_ADD_API_TOKEN_PATH is the path to handle add api token requests.
	DEVELOPER_DROP_API_TOKEN_PATH string = "/dev/api_token/drop" //DEVELOPER_DROP_API_TOKEN_PATH is the path to handle drop api token requests.
	DEVELOPER_TOKEN_LIST_PATH     string = "/dev/api_token/list" //DEVELOPER_LIST_API_TOKEN_PATH is the path to handle list api token requests.

	GAME_OWNER_AUTH_PATH         string = "/owner/auth"        //GAME_OWNER_AUTH_PATH is the path to handle GAME_OWNER authentication.
	GAME_OWNER_REGISTRATION_PATH string = "/owner/register"    //GAME_OWNER_REGISTRATION_PATH is the path to handle GAME_OWNER registration.
	GAME_OWNER_ADD_GAME_PATH     string = "/owner/game/add"    //GAME_OWNER_ADD_GAME_PATH is the path to handle add game requests.
	GAME_OWNER_REMOVE_GAME_PATH  string = "/owner/game/remove" //GAME_OWNER_REMOVE_GAME_PATH is the path to handle remove game requests.
	GAME_OWNER_GAME_LIST_PATH    string = "/owner/game/list"   //GAME_OWNER_GAME_LIST is the path to handle requests to show games of a GAME_OWNER.

	GAME_ENABLE_DISABLE_PATH string = "/game/action" //GAME_ENABLE_DISABLE_PATH is the path to handle game actions for game owner (enable, disable for users)

	//misc
	MAX_NUMBER_SALT        int           = 20000            //MAX_NUMBER_SALT is the base salt used in password hashing.
	INVALID_TOKEN          string        = "INVALID"        //INVALID_TOKEN represents an invalid token returned from a func with errors during the creation.
	DEBUG                  bool          = true             //DEBUG if true, means that application is being debugged.
	CACHE_REFRESH_INTERVAL time.Duration = time.Minute * 30 //CACHE_REFRESH_INTERVAL is the the time between cache refreshes.

	//sets
	LOGGED_USERS_SET      string = "users"             //LOGGED_USERS_SET represents the name of session set in cache of logged users.
	LOGGED_DEVELOPERS_SET string = "developers"        //LOGGED_DEVELOPERS_SET represents the name of session set in cache of logged developers.
	LOGGED_OWNERS_SET     string = "owners"            //LOGGED_OWNERS_SET represents the name of session set in cache of logged owners.
	API_TOKENS_SET        string = "API_Tokens"        //API_TOKENS_SET represents the name of set in cache of most used API tokens.
	SUMMARY_GAMES_SET     string = "games/summary/all" //SUMMARY_GAMES_SET represents the name of the set of all summaries of the games (name + desc only).

	/* entities */
	USER_IN_CACHE                 string = LOGGED_USERS_SET + "/with_id/%d"      //USER_IN_CACHE represents how to find a single user from the ID in the cache.
	USER_ENABLED_GAMES            string = USER_IN_CACHE + "/enabled_games"      //USER_ENABLED_GAMES represents how to find enabled games for a user in cache.
	OWNER_IN_CACHE                string = LOGGED_OWNERS_SET + "/with_id/%d"     //OWNER_IN_CACHE represents how to find a single owner from the ID in the cache.
	DEVELOPER_IN_CACHE            string = LOGGED_DEVELOPERS_SET + "/with_id/%d" //DEVELOPER_IN_CACHE represents how to find a single developer from the ID in the cache.
	DEVELOPER_TOKEN_LIST_IN_CACHE string = DEVELOPER_IN_CACHE + "/token_list"    //DEVELOPER_IN_CACHE represents how to find the token list of a developer in cache.
)
