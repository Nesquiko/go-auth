package mocks

import (
	"database/sql"
	"fmt"

	"github.com/Nesquiko/go-auth/pkg/db"
	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

type DBConnectionMock struct {
}

var fakeDB = make(map[string]*db.UserDBEntity)

func (dbConn DBConnectionMock) UserByUsername(username string) (*db.UserDBEntity, error) {
	user, ok := fakeDB[username]
	if !ok {
		return nil, sql.ErrNoRows
	}

	return user, nil
}

func (dbConn DBConnectionMock) SaveUser(user *db.UserModel) error {

	if _, ok := fakeDB[user.Username]; ok {
		return &mysql.MySQLError{
			Number:  1062,
			Message: fmt.Sprintf("duplicate entry '%s' for users.username", user.Username),
		}
	}

	for _, v := range fakeDB {
		if v.Email == user.Email {
			return &mysql.MySQLError{
				Number:  1062,
				Message: fmt.Sprintf("duplicate entry '%s' for users.email", user.Email),
			}
		}
	}

	fakeDB[user.Username] = &db.UserDBEntity{
		Uuid:         uuid.New(),
		Email:        user.Email,
		Username:     user.Username,
		PasswordHash: user.PasswordHash,
	}

	return nil
}

func (dbConn DBConnectionMock) Save2FASecret(username, secret string) error {
	fakeDB[username].Secret2FA = sql.NullString{String: secret}

	return nil
}
func (dbConn DBConnectionMock) Get2FASecret(username string) (string, error) {
	panic("Not implemented")
}

func (dbConn DBConnectionMock) UpdateEnabled2FA(username string, enabled bool) error {
	panic("Not implemented")
}

func (dbConn DBConnectionMock) GetEnabled2FA(username string) (bool, error) {
	panic("Not implemented")
}
