package db

import (
	"fmt"
)

// UserByUsername returns a UserDBEntity from database specified by the username
// parameter. If the username doesn't exist, sql.ErrNoRows error is returned.
func (db connection) UserByUsername(username string) (*UserDBEntity, error) {
	var user UserDBEntity
	var enabled2FAStr string

	row := db.QueryRow("SELECT * FROM users WHERE username = ?", username)

	if err := row.Scan(&user.Uuid, &user.Username, &user.Email, &user.PasswordHash,
		&user.Secret2FA, &enabled2FAStr); err != nil {
		fmt.Println(err)
		return nil, err
	}

	if enabled2FAStr == "\x00" {
		user.Enabled2FA = false
	} else {
		user.Enabled2FA = true
	}

	return &user, nil
}

// SaveUser saves the UserModel passed as parameter to a database.
func (db connection) SaveUser(user *UserModel) error {
	_, err := db.Exec(
		"INSERT INTO users (username, email, passwordHash, secret2FA, enabled2FA) VALUES (?, ?, ?, ?, ?)",
		user.Username,
		user.Email,
		user.PasswordHash,
		nil,
		0,
	)

	if err != nil {
		return err
	}

	return nil
}

// Save2FASecret saves new secret needed during 2FA
func (db connection) Save2FASecret(username, secret string) error {
	_, err := db.Exec(
		"UPDATE users SET secret2FA = ? WHERE username = ?",
		secret,
		username,
	)

	if err != nil {
		return err
	}

	return nil
}
