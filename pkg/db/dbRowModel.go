package db

import (
	"fmt"

	"github.com/google/uuid"
)

type UserDBEntity struct {
	Uuid         uuid.UUID
	Username     string
	Email        string
	PasswordHash string
}

func (u UserDBEntity) String() string {
	return fmt.Sprintf("username: %s | email: %s | uuid: %s",
		u.Username,
		u.Email,
		u.Uuid.String())
}

type UserModel struct {
	Username     string
	Email        string
	PasswordHash string
}

func (u UserModel) String() string {
	return fmt.Sprintf("username: %s | email: %s",
		u.Username,
		u.Email)
}
