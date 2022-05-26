package db

// UserByUsername returns a UserDBEntity from database specified by the username
// parameter. If the username doesn't exist, sql.ErrNoRows error is returned.
func (db connection) UserByUsername(username string) (*UserDBEntity, error) {
	var user UserDBEntity

	row := db.QueryRow("SELECT * FROM users WHERE username = ?", username)

	if err := row.Scan(&user.Uuid, &user.Username, &user.Email, &user.PasswordHash); err != nil {
		return nil, err
	}
	return &user, nil
}

// SaveUser saves the UserModel passed as parameter to a database.
func (db connection) SaveUser(user *UserModel) error {
	_, err := db.Exec(
		"INSERT INTO users (username, email, passwordHash) VALUES (?, ?, ?)",
		user.Username,
		user.Email,
		user.PasswordHash,
	)

	if err != nil {
		return err
	}

	return nil
}
