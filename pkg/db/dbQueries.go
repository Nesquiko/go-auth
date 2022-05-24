package db

func (db connection) UserByUsername(username string) (*UserDBEntity, error) {
	var user UserDBEntity

	row := db.QueryRow("SELECT * FROM users WHERE username = ?", username)

	if err := row.Scan(&user.Uuid, &user.Username, &user.Email, &user.PasswordHash); err != nil {
		return nil, err
	}
	return &user, nil
}

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
