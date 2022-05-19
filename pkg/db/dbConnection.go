package db

import (
	"database/sql"

	"github.com/go-sql-driver/mysql"
)

type Connection struct {
	*sql.DB
}

func Connect(user, passwd, dbname, addr string) (*Connection, error) {
	var db *sql.DB
	cfg := mysql.Config{
		User:   user,
		Passwd: passwd,
		Net:    "tcp",
		Addr:   addr,
		DBName: dbname,
	}

	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())

	if err != nil {
		return nil, err
	}

	pingErr := db.Ping()
	if pingErr != nil {
		return nil, pingErr
	}

	return &Connection{db}, nil
}
