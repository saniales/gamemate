package test

import (
    "testing"
    "database/sql"
    "fmt"
  _ "github.com/go-sql-driver/mysql"
)

const(
  TEST_USER = "test_gamemate"
  TEST_PASSWORD = "test"
)

func TestLinkMySQL(test *testing.T) {
  if testing.Short() {
    test.Skip("Skipping connection Test in short mode.")
  } else {
    var test_handle *sql.DB
    var err error
    test_handle, err = configurations.LinkMySQL(TEST_USER, TEST_PASSWORD);
    if err != nil {
      test.Log("Error creating DB Object : error => " + err.Error())
      test.Fail()
    } else {
      err = test_handle.Ping()
      if err != nil {
        test.Log("Error while pinging the DB : Cannot perform PING. Error => " + err.Error())
        test.Fail()
      } else {
        //try query SELECT * FROM test_table and check some values
      }
    }
  }
}
