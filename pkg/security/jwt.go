package security

import (
	"crypto/rand"
	"time"

	"github.com/golang-jwt/jwt"
)

// jwtKey is used as a secret when generating and validating a JWT token.
var jwtKey []byte

// expirationDuration is how long a JWT token is valid.
var expirationDuration time.Duration = 5 * time.Minute

// init function creates new key using crypto/rand. If the creation of the key
// is not successful, it panics.
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

// claims represents JWT claims used in body of JWT.
type claims struct {
	Username      string `json:"username"`
	Authenticated bool   `json:"authenticated"`
	jwt.StandardClaims
}

// GenerateUnauthenticatedJWT generates new JWT with a username as a claim and
// authenticated claim set to false, because 2FA is needed to be fully
// authenticated. The JWT has an expiration time equal to the expirationDuration
// variable. The signing algorithm is HS256.
func GenerateUnauthenticatedJWT(username string) (string, error) {

	expirationTime := time.Now().Add(expirationDuration)

	claims := &claims{
		Username:      username,
		Authenticated: false,
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
