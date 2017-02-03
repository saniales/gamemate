package libs

import "github.com/garyburd/redigo/redis" //Package to interact with Redis DB

//Represents a logged user into the system.
type User struct {
    Username string `json:"Username" xml:"Username" form:"Username"`
    Email    string `json:"Email" xml:"Email" form:"Email"`
    Birthday int64  `json:"Birthday" xml:"Birthday" form:"Birthday"`
    Gender   string `json:"Gender" xml:"Gender" form:"Gender"`
    //password -> a SHA_512 Hash of the password (not salted - on DB - not in this struct).
    //salt -> a number to enhance hash (on DB - not in this struct).
}

func (u User) InsertIntoCache(pool *redis.Pool) error {
    var conn redis.Conn = pool.Get()
    defer conn.Close()

    _, err := conn.Do("HMSET", "user:" + u.Username, "Username", u.Username, "Email", u.Email, "Birthday", u.Birthday, "Gender", u.Gender)
    return err;
}

func (u User) DeleteFromCache(pool *redis.Pool) error {
  conn := pool.Get()
  defer conn.Close()

  _, err := conn.Do("DEL", "user:" + u.Username)
  return err;
}
