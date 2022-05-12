package server

import (
	"golang.org/x/crypto/bcrypt"
)

func encryptPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

// Comparing the password with the hash
// err = bcrypt.CompareHashAndPassword(hashedPassword, password)
// fmt.Println(err) // nil means it is a match
