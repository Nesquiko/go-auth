package db

import (
	"database/sql"
	"fmt"
	"log"
)

func (db connection) FetchUsers() ([]UserDBEntity, error) {
	var users []UserDBEntity

	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		return nil, fmt.Errorf("error from query: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var user UserDBEntity
		if err := rows.Scan(&user.Uuid, &user.Username, &user.Email, &user.PasswordHash); err != nil {
			return nil, fmt.Errorf("error in scanning: %v", err)
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error: %v", err)
	}
	return users, nil
}

func (db connection) UserByUsername(username string) (*UserDBEntity, error) {

	var user UserDBEntity

	row := db.QueryRow("SELECT * FROM users WHERE username = ?", username)

	if err := row.Scan(&user.Uuid, &user.Username, &user.Email, &user.PasswordHash); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user by %s: %v", username, err)
		}
		return nil, fmt.Errorf("user by %s: %v", username, err)
	}
	return &user, nil
}

func (db connection) SaveUser(user *UserModel) error {
	result, err := db.Exec(
		"INSERT INTO users (username, email, passwordHash) VALUES (?, ?, ?)",
		user.Username,
		user.Email,
		user.PasswordHash,
	)
	if err != nil {
		return fmt.Errorf("saveUser %s: %v", user.Username, err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("saveUser %s: %v", user.Username, err)
	}
	log.Println(id)

	return nil
}
