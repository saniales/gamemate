package test

import (
	"fmt"
	"sanino/gamemate/configurations"
	"testing"
)

func TestInitServer(test *testing.T) {
	var test_server = configurations.InitServer()
	if test_server == nil {
		test.Log("test server not initialized correctly, there must be a problem with the default parameters.")
		test.Fail()
	} else {
		fmt.Println("OK")
		//Output: OK
	}
}
