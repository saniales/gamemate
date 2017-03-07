package test

import (
	"sanino/gamemate/configurations"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

const (
	TEST_USER     = "test_gamemate"
	TEST_PASSWORD = "test"
)

func TestLinkMySQL(test *testing.T) {
	if testing.Short() {
		test.Skip("Skipping connection Test in short mode.")
	} else {
		var err error
		err = configurations.InitArchivesWithAuth(TEST_USER, TEST_PASSWORD)
		test_handle := configurations.ArchivesPool
		if err != nil {
			test.Log("Error creating DB Object : error => " + err.Error())
			test.Fail()
		} else {
			err = test_handle.Ping()
			if err != nil {
				test.Log("Error while pinging the DB : Cannot perform PING. Error => " + err.Error())
				test.Fail()
			} else {
				test.Log(test_handle.Query("SELECT * FROM test_table"))
			}
		}
	}
}
