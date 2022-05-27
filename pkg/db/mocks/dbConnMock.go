package mocks

import (
	"database/sql"

	"github.com/Nesquiko/go-auth/pkg/db"
	"github.com/go-sql-driver/mysql"
)

type DBConnectionMock struct {
}

var fakeDB = make(map[string]*db.UserDBEntity, 0)

func (db DBConnectionMock) UserByUsername(username string) (*db.UserDBEntity, error) {
	user, ok := fakeDB[username]
	if !ok {
		return nil, sql.ErrNoRows
	}

	return user, nil
}

func (db DBConnectionMock) SaveUser(user *db.UserModel) error {

	if _, ok := fakeDB[user.Username]; ok {
		return &mysql.MySQLError{
			Number:  1062,
			Message: "duplicate entries for users.username",
		}
	}

	for _, v := range fakeDB {
		if v.Email == user.Email {
			return &mysql.MySQLError{
				Number:  1062,
				Message: "duplicate entries for users.email",
			}
		}
	}

	return nil
}
