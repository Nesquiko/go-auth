package server

import (
	"testing"
)

func Test_encryptPasswordSamePasswordsHashDoNotMatch(t *testing.T) {
	passwd := "123"
	hash1, _ := encryptPassword(passwd)
	hash2, _ := encryptPassword(passwd)

	if hash1 == hash2 {
		t.Fatalf("Hashes of same password match, they should not")
	}
}

func Test_compareHashAndPasswordValid(t *testing.T) {
	passwd := "123"
	hash1, _ := encryptPassword(passwd)

	isValid := compareHashAndPassword(hash1, passwd)
	if !isValid {
		t.Fatalf("Comparison failed, but expected not to")
	}
}

func Test_compareHashAndPasswordInvalid(t *testing.T) {
	passwd := "123"
	hash1, _ := encryptPassword(passwd)
	passwd2 := "invalid"

	isValid := compareHashAndPassword(hash1, passwd2)
	if isValid {
		t.Fatalf("Comparison succeded, but expected not to")
	}
}
