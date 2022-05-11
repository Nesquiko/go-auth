package db

import "fmt"

func (db Connection) FetchUsers() ([]User, error) {
	var users []User

	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		return nil, fmt.Errorf("error from query: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var user User
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
