package constants

import "time"

//Paths Constants
const (
	ROOT_PATH              string = "/"              //Path to the root directory of the server.
	USER_REGISTRATION_PATH string = "/user/register" //Path to handle user registration.
	AUTH_PATH              string = "/user/auth"     //Path to handle user authentication.
	GET_USER_REQUEST_PATH  string = "/user/info"     //Gets info about a user.

	ROOM_CREATION_PATH    string = "/room/create" //Path to handle the creation of a new room.
	GET_ROOM_REQUEST_PATH string = "/room/get"    //Path to handle the request of data of a particular room.

	MATCH_CREATION_PATH     string = "/match/create" //Path to handle the creationd of a match (not started).
	MATCH_START_PATH        string = "/match/start"  //Path to handle the start of a match (it becomes LIVE).
	MATCH_DATA_REQUEST_PATH string = "/match/get"    //Path to handle the request of getting data of a particular match.

	TURN_ACTION_PATH string = "/user/move/make"     //Path to handle an action in a match.
	TURN_END_PATH    string = "/user/move/end_turn" //Path to handle the end of a turn.

	DEVELOPER_AUTH_PATH          string = "/dev/auth"           //Path to handle developer authentication.
	DEVELOPER_REGISTRATION_PATH  string = "/dev/register"       //Path to hande developer registration.
	DEVELOPER_ADD_API_TOKEN_PATH string = "/dev/api_token/add"  //Path to handle add api token requests.
	DEVELOPER_DROP_API_TOKEN     string = "/dev/api_token/drop" //Path to handle drop api token requests.

	VENDOR_AUTH_PATH         string = "/vendor/auth"        //Path to handle vendor authentication.
	VENDOR_REGISTRATION_PATH string = "/vendor/register"    //Path to handle vendor registration.
	VENDOR_ADD_GAME_PATH     string = "/vendor/game/add"    //Path to handle add game requests.
	VENDOR_REMOVE_GAME_PATH  string = "/vendor/game/remove" //Path to handle remove game requests.
	VENDOR_GAME_LIST         string = "/vendor/game/list"   //Path to handle requests to show games of a vendor.

	//misc
	MAX_NUMBER_SALT        int           = 20000               //base salt used in password hashing.
	INVALID_TOKEN          string        = "INVALID"           //Represents an invalid token returned from a func with errors during the creation.
	CACHE_REFRESH_INTERVAL time.Duration = time.Minute * 30    //The time between cache refreshes.
	LOGGED_USERS_SET       string        = "logged_users"      //represents the name of session set in cache of logged users.
	LOGGED_DEVELOPERS_SET  string        = "logged_developers" //represents the name of session set in cache of logged developers.
	DEBUG                  bool          = true                //if true, application is being debugged.
)
