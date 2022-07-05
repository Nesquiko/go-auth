package db

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
)

func TestUserDBEntity_String(t *testing.T) {
	id, _ := uuid.NewUUID()
	username, email, passwdHash := "Johny", "john@foo.bar", "$asdfioj5641684"

	want := fmt.Sprintf(
		"username: %s | email: %s | enabled 2FA: %v | uuid: %s",
		username,
		email,
		false,
		id.String(),
	)

	u := UserDBEntity{Uuid: id, Username: username, Email: email, PasswordHash: passwdHash}

	if got := u.String(); got != want {
		t.Errorf("UserDBEntity.String() = %v, want %v", got, want)
	}
}

func TestUserModel_String(t *testing.T) {
	username, email, passwdHash := "Johny", "john@foo.bar", "$asdfioj5641684"
	want := fmt.Sprintf("username: %s | email: %s", username, email)

	u := UserModel{Username: username, Email: email, PasswordHash: passwdHash}
	if got := u.String(); got != want {
		t.Errorf("UserModel.String() = %v, want %v", got, want)
	}
}
