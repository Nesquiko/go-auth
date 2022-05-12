package db

import (
	"database/sql"

	"github.com/go-sql-driver/mysql"
)

type Connection struct {
	*sql.DB
}

func Connect(user, passwd, dbname, addr string) (*Connection, func(), error) {
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
		return nil, nil, err
	}

	pingErr := db.Ping()
	if pingErr != nil {
		return nil, nil, err
	}

	return &Connection{db}, func() { db.Close() }, err
}
