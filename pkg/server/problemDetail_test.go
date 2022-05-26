package server

import (
	"database/sql"
	"errors"
	"net/http"
	"reflect"
	"testing"

	"github.com/go-sql-driver/mysql"
)

func TestGetProblemDetailsMySQLDuplicateEntryUsername(t *testing.T) {
	fakeError := mysql.MySQLError{
		Number:  1062,
		Message: "duplicate entries for users.username",
	}
	wantCode := http.StatusConflict
	wantType := "username.already_exists"

	problem := GetProblemDetails(&fakeError)

	if problem.StatusCode != wantCode {
		t.Errorf("Wrong status code, expected: %d, but was %d", wantCode, problem.StatusCode)
	}
	if problem.Type != wantType {
		t.Errorf("Wrong type, expected: %q, but was %q", wantType, problem.Type)
	}
}

func TestGetProblemDetailsMySQLDuplicateEntryEmail(t *testing.T) {
	fakeError := mysql.MySQLError{
		Number:  1062,
		Message: "duplicate entries for users.email",
	}
	wantCode := http.StatusConflict
	wantType := "email.already_used"

	problem := GetProblemDetails(&fakeError)

	if problem.StatusCode != wantCode {
		t.Errorf("Wrong status code, expected: %d, but was %d", wantCode, problem.StatusCode)
	}
	if problem.Type != wantType {
		t.Errorf("Wrong type, expected: %q, but was %q", wantType, problem.Type)
	}
}

func TestGetProblemDetailsMySQLDuplicateUnknownEntry(t *testing.T) {
	fakeError := mysql.MySQLError{
		Number:  1062,
		Message: "duplicate entries",
	}
	wantCode := http.StatusConflict
	wantType := "unknown.duplicate"

	problem := GetProblemDetails(&fakeError)

	if problem.StatusCode != wantCode {
		t.Errorf("Wrong status code, expected: %d, but was %d", wantCode, problem.StatusCode)
	}
	if problem.Type != wantType {
		t.Errorf("Wrong type, expected: %q, but was %q", wantType, problem.Type)
	}
}

func TestGetProblemDetailsUnknownMySQLError(t *testing.T) {
	fakeError := mysql.MySQLError{
		Number:  1,
		Message: "unknown error",
	}

	problem := GetProblemDetails(&fakeError)

	if !reflect.DeepEqual(problem, UnexpectedErrorProblem) {
		t.Errorf("Returned problem is not UnexpectedErrorProblem, %q", problem)
	}
}

func TestGetProblemDetailsSQLNoRowsError(t *testing.T) {
	fakeError := sql.ErrNoRows
	wantCode := http.StatusUnauthorized
	wantType := "username.not_found"

	problem := GetProblemDetails(fakeError)

	if problem.StatusCode != wantCode {
		t.Errorf("Wrong status code, expected: %d, but was %d", wantCode, problem.StatusCode)
	}
	if problem.Type != wantType {
		t.Errorf("Wrong type, expected: %q, but was %q", wantType, problem.Type)
	}
}

func TestGetProblemDetailsUnknownError(t *testing.T) {
	fakeError := errors.New("unknown error")

	problem := GetProblemDetails(fakeError)

	if !reflect.DeepEqual(problem, UnexpectedErrorProblem) {
		t.Errorf("Returned problem is not UnexpectedErrorProblem, %q", problem)
	}
}
