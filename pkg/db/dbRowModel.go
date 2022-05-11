package db

import "fmt"

type User struct {
	Uuid         string
	Username     string
	Email        string
	PasswordHash string
}

func (u User) String() string {
	return fmt.Sprintf("username: %s | email: %s | uuid: %s",
		u.Username,
		u.Email,
		u.Uuid)
}
