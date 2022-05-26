// Package security provides functions for encryption of passwords
// and JWT token handling.
package security

import (
	"golang.org/x/crypto/bcrypt"
)

// EncryptPassword encrypts the given password with Bcrypt algorithm and returns
// the generated hash.
func EncryptPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

// HashAndPasswordMatch takes a hash and a password and determines if the hash
// matches the entered password.
func HashAndPasswordMatch(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
