package constants

import "time"

//Paths Constants
const (
	ROOT_PATH              string = "/"              //Path to the root directory of the server.
	USER_REGISTRATION_PATH string = "/register_user" //Path to handle user registration.
	AUTH_PATH              string = "/auth_user"     //Path to handle user authentication.
	GET_USER_REQUEST_PATH  string = "/user_info"     //Gets info about a user.
	ROOM_CREATION_PATH     string = "/new_room"      //Path to handle the creation of a new room (only with thde
	//user who requested the creation). The new room will be open.
	GET_ROOM_REQUEST_PATH   string = "/get_room"      //Path to handle the request of data of a particular rodom.
	MATCH_CREATION_PATH     string = "/new_match"     //Path to handle the creationd of a match (not started).
	MATCH_START_PATH        string = "/start_match"   //Path to handle the start of a match (it becomes LIVE).
	MATCH_DATA_REQUEST_PATH string = "/get_match"     //Path to handle the request of getting data of a particular match.
	TURN_ACTION_PATH        string = "/make_move"     //Path to handle an action in a match.
	TURN_END_PATH           string = "/end_turn"      //Path to handle the end of a turn.
	TURN_END_MATCH_ACK      string = "/end_match_ack" //Path to let the server know that the client received the

	MAX_NUMBER_SALT        int           = 20000               //base salt used in password hashing.
	INVALID_TOKEN          string        = "INVALID"           //Represents an invalid token returned from a func with errors during the creation.
	CACHE_REFRESH_INTERVAL time.Duration = time.Minute * 30    //The time between cache refreshes.
	LOGGED_USERS_SET       string        = "logged_users"      //represents the name of session set in cache of logged users.
	LOGGED_DEVELOPERS_SET  string        = "logged_developers" //represents the name of session set in cache of logged developers.
	DEBUG                  bool          = true                //if true, application is being debugged.
)
