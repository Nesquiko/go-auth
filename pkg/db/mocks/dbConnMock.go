package mocks

import "github.com/Nesquiko/go-auth/pkg/db"

type DBConnectionMock struct {
}

func (db DBConnectionMock) UserByUsername(username string) (*db.UserDBEntity, error) {

	return nil, nil
}

func (db DBConnectionMock) SaveUser(user *db.UserModel) error {
	return nil
}
