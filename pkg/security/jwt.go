package security

import (
	"crypto/rand"
	"time"

	"github.com/golang-jwt/jwt"
)

var jwtKey []byte
var expirationDuration time.Duration = 5 * time.Hour

func init() {
	keyLen := 64
	jwtKey = make([]byte, keyLen)
	n, err := rand.Read(jwtKey)

	if n != keyLen {
		panic("Didn't create a secret")
	}
	if err != nil {
		panic("Error occured wihle creating a secret")
	}
}

type claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func GenerateJWT(username string) (string, error) {

	expirationTime := time.Now().Add(expirationDuration)

	claims := &claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		}}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
