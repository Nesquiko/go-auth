package security

import (
	"encoding/base64"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestGenerateJWTPayloadCorrectUsernameClaim(t *testing.T) {
	username := "Joe"
	wantClaim := fmt.Sprintf("%q:%q", "username", username)
	jwt, err := GenerateJWT(username, false)

	if err != nil {
		t.Fatalf("err was not nil, %q", err.Error())
	}
	jwtSplit := strings.Split(jwt, ".")

	payload, err := base64.RawURLEncoding.DecodeString(jwtSplit[1])
	if err != nil {
		t.Fatalf("decoding err was not nil, %q", err.Error())
	}
	if !strings.Contains(string(payload), wantClaim) {
		t.Errorf("No valid username claim in payload %s", payload)
	}
}

func TestGenerateJWTPayloadCorrectExpClaim(t *testing.T) {
	username := "Joe"
	exp := time.Now().Add(expirationDuration)
	wantClaim := fmt.Sprintf("%q:%d", "exp", exp.Unix())
	jwt, err := GenerateJWT(username, false)

	if err != nil {
		t.Fatalf("err was not nil, %q", err.Error())
	}
	jwtSplit := strings.Split(jwt, ".")

	payload, err := base64.RawURLEncoding.DecodeString(jwtSplit[1])
	if err != nil {
		t.Fatalf("decoding err was not nil, %q", err.Error())
	}
	if !strings.Contains(string(payload), wantClaim) {
		t.Errorf("No valid exp claim in payload %s", payload)
	}
}
