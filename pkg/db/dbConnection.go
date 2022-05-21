package db

import (
	"database/sql"

	"github.com/go-sql-driver/mysql"
)

type connection struct {
	*sql.DB
}

var DBConnection *connection = nil

func ConnectDB(driver, dsn string) error {
	var err error

	if DBConnection == nil {
		DBConnection, err = connect(driver, dsn)
	}

	return err
}

func MySQLDSNConfig(user, passwd, addr, dbname string) *mysql.Config {
	return &mysql.Config{
		User:   user,
		Passwd: passwd,
		Net:    "tcp",
		Addr:   addr,
		DBName: dbname,
	}
}

func connect(driver, dsn string) (*connection, error) {
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
