package configurations

import (
    "database/sql"
    "fmt"
  _ "github.com/go-sql-driver/mysql"
)

const(
  MYSQL_DRIVER = "mysql"
  MYSQL_USER = "user"
  MYSQL_PASSWORD = "password"
  MYSQL_DB_NAME = "dbname"
  MYSQL_HOST = "localhost"
)

//Creates a MySQL communication point for the API for the Data.
func LinkMySQL(user string, password string) (*DB, error) {
    var database_handle *DB
    var err error

    //creates database name space to connect to database
    var database_namespace string = fmt.Sprintf("%s:%s@https(%s)/%s", MYSQL_USER, MYSQL_PASSWORD, MYSQL_HOST, MYSQL_DB_NAME)

    //creates a handle to the DB, keep in mind that there is a pool of connections in background
    //and a connection is not open if it isn't needed.
    //Use database_handle.Ping() to verify if a DB is connected
    database_handle, err = sql.Open(MYSQL_DRIVER, database_namespace)
    if err != nil {
      return nil, err
    } else {
      //checking if the host is reachable using the specified DSN (Database NameSpace)
      err = database_handle.Ping()
      if err != nil {
        return nil, err
      } else {
        return database_handle, nil
      }
    }
}

func LinkMySQL() (*DB, error){
  return LinkMySQL(MYSQL_USER, MYSQL_PASSWORD)
}
