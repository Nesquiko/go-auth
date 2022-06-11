package db

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

// UserDBEntity is a model for how users are represented in database.
type UserDBEntity struct {
	// Uuid is a primary key of the UserDBEntity entry in database.
	Uuid uuid.UUID

	// Username of the user, must be unique.
	Username string

	// Email of the user, must be unique.
	Email string

	// Password hash of users account password.
	PasswordHash string

	// Secret2FA is used during 2FA
	Secret2FA sql.NullString

	// Enabled2FA indicates if user enabled 2FA
	Enabled2FA bool
}

// String returns string representation of a UserDBEntity.
func (u UserDBEntity) String() string {
	return fmt.Sprintf("username: %s | email: %s | enabled 2FA: %v | uuid: %s",
		u.Username,
		u.Email,
		u.Enabled2FA,
		u.Uuid.String())
}

// UserModel represents a model of an user, which is used in application logic.
type UserModel struct {
	// Username of the user, must be unique.
	Username string

	// Email of the user, must be unique.
	Email string

	// Password hash of users account password.
	PasswordHash string
}

// String returns string representation of a UserModel.
func (u UserModel) String() string {
	return fmt.Sprintf("username: %s | email: %s",
		u.Username,
		u.Email,
	)
}
