package OutGameRequests

//GameList represents a request to all the games available to a single user.
type GameList struct {
	Type         string
	SessionToken string
}
