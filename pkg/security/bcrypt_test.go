package security

import (
	"testing"
)

func Test_encryptPasswordSamePasswordsHashDoNotMatch(t *testing.T) {
	passwd := "123"
	hash1, _ := EncryptPassword(passwd)
	hash2, _ := EncryptPassword(passwd)

	if hash1 == hash2 {
		t.Fatalf("Hashes of same password match, they should not")
	}
}

func Test_compareHashAndPasswordValid(t *testing.T) {
	passwd := "123"
	hash1, _ := EncryptPassword(passwd)

	isValid := HashAndPasswordMatch(hash1, passwd)
	if !isValid {
		t.Fatalf("Comparison failed, but expected not to")
	}
}

func Test_compareHashAndPasswordInvalid(t *testing.T) {
	passwd := "123"
	hash1, _ := EncryptPassword(passwd)
	passwd2 := "invalid"

	isValid := HashAndPasswordMatch(hash1, passwd2)
	if isValid {
		t.Fatalf("Comparison succeded, but expected not to")
	}
}
