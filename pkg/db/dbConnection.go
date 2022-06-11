// Package db provides functions for interacting with a database. Only MySQL
// database is supported for now.
package db

import (
	"database/sql"

	"github.com/go-sql-driver/mysql"
)

// DBConnection represents a layer between the database and the application
// logic.
type DBConnection interface {

	// UserByUsername returns a UserDBEntity from database specified by the username
	// parameter. If the username doesn't exist, error is returned.
	UserByUsername(username string) (*UserDBEntity, error)

	// SaveUser saves the UserModel passed as parameter to a database.
	SaveUser(user *UserModel) error

	// Save2FASecret saves secret for 2FA
	Save2FASecret(username, secret string) error
}

// connection struct with embedded sql.DB struct serving as a layer between
// application logic and database logic. Also this struct implements the
// DBConnection interface.
type connection struct {
	*sql.DB
}

// DBConnectino is a global connection to a database, through which application
// logic interacts with database.
var DBConn DBConnection = nil

// ConnectDB establishes the global connection to a database.
func ConnectDB(driver, dsn string) error {
	var err error

	if DBConn == nil {
		DBConn, err = connect(driver, dsn)
	}

	return err
}

// MySQLDSNConfig is a simple util function for creating a mysql.Config with
// custom connection options.
func MySQLDSNConfig(user, passwd, addr, dbname string) *mysql.Config {
	return &mysql.Config{
		User:   user,
		Passwd: passwd,
		Net:    "tcp",
		Addr:   addr,
		DBName: dbname,
	}
}

// connect tries to establish a connection to a database. The param driver
// specifies the type of the database and dsn is the configuration used to
// connect to the database. After initializing the connection to the database,
// it is pinged to see if the connection was established.
func connect(driver, dsn string) (DBConnection, error) {
	var db *sql.DB

	db, err := sql.Open(driver, dsn)

	if err != nil {
		return nil, err
	}

	pingErr := db.Ping()
	if pingErr != nil {
		return nil, pingErr
	}

	return &connection{db}, nil
}
