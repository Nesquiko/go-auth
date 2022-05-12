package db

import (
	"fmt"

	"github.com/google/uuid"
)

type User struct {
	Uuid         uuid.UUID
	Username     string
	Email        string
	PasswordHash string
}

func (u User) String() string {
	return fmt.Sprintf("username: %s | email: %s | uuid: %s",
		u.Username,
		u.Email,
		u.Uuid.String())
}
