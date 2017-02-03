package test

import (
  "testing"
  "sanino/gamemate/configurations"
  "github.com/labstack/echo"
  "fmt"
)

func TestInitServer(test *testing.T)  {
  var test_server *echo.Echo = configurations.InitServer()
  if test_server == nil {
    test.Log("test server not initialized correctly, there must be a problem with the default parameters.")
    test.Fail()
  } else {
    fmt.Println("OK")
    //Output: OK
  }
}
