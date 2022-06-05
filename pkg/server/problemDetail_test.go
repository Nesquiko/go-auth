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
	wantInstance := "/signup"

	problem := GetProblemDetails(&fakeError, wantInstance)

	if problem.StatusCode != wantCode {
		t.Errorf("Wrong status code, expected: %d, but was %d", wantCode, problem.StatusCode)
	}
	if problem.Instance != wantInstance {
		t.Errorf("Wrong instance, expected: %q, but was %q", wantInstance, problem.Instance)
	}
}

func TestGetProblemDetailsMySQLDuplicateEntryEmail(t *testing.T) {
	fakeError := mysql.MySQLError{
		Number:  1062,
		Message: "duplicate entries for users.email",
	}
	wantCode := http.StatusConflict
	wantInstance := "/signup"

	problem := GetProblemDetails(&fakeError, wantInstance)

	if problem.StatusCode != wantCode {
		t.Errorf("Wrong status code, expected: %d, but was %d", wantCode, problem.StatusCode)
	}
	if problem.Instance != wantInstance {
		t.Errorf("Wrong instance, expected: %q, but was %q", wantInstance, problem.Instance)
	}
}

func TestGetProblemDetailsMySQLDuplicateUnknownEntry(t *testing.T) {
	fakeError := mysql.MySQLError{
		Number:  1062,
		Message: "duplicate entries",
	}
	wantCode := http.StatusConflict
	wantInstance := "/signup"

	problem := GetProblemDetails(&fakeError, wantInstance)

	if problem.StatusCode != wantCode {
		t.Errorf("Wrong status code, expected: %d, but was %d", wantCode, problem.StatusCode)
	}
	if problem.Instance != wantInstance {
		t.Errorf("Wrong instance, expected: %q, but was %q", wantInstance, problem.Instance)
	}
}

func TestGetProblemDetailsUnknownMySQLError(t *testing.T) {
	fakeError := mysql.MySQLError{
		Number:  1,
		Message: "unknown error",
	}

	wantInstance := "/login"
	problem := GetProblemDetails(&fakeError, wantInstance)

	if !reflect.DeepEqual(problem, UnexpectedErrorProblem(wantInstance)) {
		t.Errorf("Returned problem is not UnexpectedErrorProblem, %q", problem)
	}
}

func TestGetProblemDetailsSQLNoRowsError(t *testing.T) {
	fakeError := sql.ErrNoRows
	wantCode := http.StatusUnauthorized
	wantInstance := "/login"

	problem := GetProblemDetails(fakeError, wantInstance)

	if problem.StatusCode != wantCode {
		t.Errorf("Wrong status code, expected: %d, but was %d", wantCode, problem.StatusCode)
	}
	if problem.Instance != wantInstance {
		t.Errorf("Wrong instance, expected: %q, but was %q", wantInstance, problem.Instance)
	}
}

func TestGetProblemDetailsUnknownError(t *testing.T) {
	fakeError := errors.New("unknown error")
	wantInstance := "/login"

	problem := GetProblemDetails(fakeError, wantInstance)

	if !reflect.DeepEqual(problem, UnexpectedErrorProblem(wantInstance)) {
		t.Errorf("Returned problem is not UnexpectedErrorProblem, %q", problem)
	}
}
