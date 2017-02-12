package configurations

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql" // used by database.sql and not by the application
)

const (
	MYSQL_DRIVER   = "mysql"
	MYSQL_USER     = "gamemate_agent"
	MYSQL_PASSWORD = "elGq9WjXWfuqqJIDP2Zu"
	MYSQL_DB_NAME  = "gamemate_archives"
	MYSQL_HOST     = "localhost:3306"
)

//ArchivesPool is a pool of connections to the archive of the system (in this case using MySQL).
var ArchivesPool *sql.DB

//ArchivesInitialized is true if a pool of connections to the archives has been initialized at least once.
var ArchivesInitialized = false

//InitArchivesWithAuth Creates a MySQL communication point for the API for the Data.
func InitArchivesWithAuth(user string, password string) error {
	var err error
	//creates database name space to connect to database
	var database_namespace string = fmt.Sprintf("%s:%s@tcp(%s)/%s", MYSQL_USER, MYSQL_PASSWORD, MYSQL_HOST, MYSQL_DB_NAME)

	//creates a handle to the DB, keep in mind that there is a pool of connections in background
	//and a connection is not open if it isn't needed.
	//Use database_handle.Ping() to verify if a DB is connected
	ArchivesPool, err = sql.Open(MYSQL_DRIVER, database_namespace)
	if err != nil {
		ArchivesInitialized = false
		return err
	}
	//checking if the host is reachable using the specified DSN (Database NameSpace)
	err = ArchivesPool.Ping()
	if err != nil {
		ArchivesInitialized = false
		return err
	}
	ArchivesInitialized = true
	return nil
}

//InitArchives links the archives with the system with the global options.
func InitArchives() error {
	return InitArchivesWithAuth(MYSQL_USER, MYSQL_PASSWORD)
}
