package test

import (
	"fmt"
	"sanino/gamemate/configurations"
	"testing"

	"github.com/garyburd/redigo/redis" //Package to interact with Redis DB
)

func TestLinkRedis(test *testing.T) {
	if testing.Short() {
		test.Skip("Skipping connection Test in short mode.")
	} else {
		var test_pool *redis.Pool = configurations.LinkRedis()
		if test_pool != nil {
			defer test_pool.Close()
			var test_connection redis.Conn = test_pool.Get()
			defer test_connection.Close()
			if test_connection != nil {
				result, err := test_connection.Do("PING")
				if err == nil {
					fmt.Println(result)
					//Output: PONG
				} else {
					test.Log("PING Command not successful : error => " + err.Error())
					test.Fail()
				}
			} else {
				test.Log("Pool not initialized : the test connection is nil")
				test.Fail()
			}
			defer test_connection.Close()
		} else {
			test.Log("Pool not initialized : the test pool is nil")
			test.Fail()
		}
	}
}
