package db

import (
	"errors"
	"log"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

var stubDB *connection
var mock sqlmock.Sqlmock
var model UserModel = UserModel{
	Username:     "James",
	Email:        "jam@bar.com",
	PasswordHash: "as46984asdfkjSDFas",
}

func TestMain(m *testing.M) {
	stubDB, mock = newMock()
	code := m.Run()
	stubDB.Close()
	os.Exit(code)
}

func Test_connectionSaveUser(t *testing.T) {
	query := "INSERT INTO users"

	mock.ExpectExec(query).
		WithArgs(model.Username, model.Email, model.PasswordHash, nil, 0).
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err := stubDB.SaveUser(&model); err != nil {
		t.Errorf("error was not expected: %s", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func Test_connectionSaveUserError(t *testing.T) {
	query := "INSERT INTO users"

	mock.ExpectExec(query).
		WithArgs(model.Username, model.Email, model.PasswordHash, nil, 0).
		WillReturnError(errors.New("testing error"))

	if err := stubDB.SaveUser(&model); err == nil {
		t.Errorf("error was expected: %s", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSave2FASecretNullWhenNotSet(t *testing.T) {
	query := "UPDATE users SET secret2FA"
	secret := "ZSOOSQWFTYYO7VZI"

	mock.ExpectExec(query).
		WithArgs(secret, model.Username).
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err := stubDB.Save2FASecret(model.Username, secret); err != nil {
		t.Errorf("error was not expected: %s", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func newMock() (*connection, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return &connection{db}, mock
}
