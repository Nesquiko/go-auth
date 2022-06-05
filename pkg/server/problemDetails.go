package server

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/Nesquiko/go-auth/pkg/api"
	"github.com/go-sql-driver/mysql"
)

// UnexpectedErrorProblem returns generic problem details response used when an
// error, about which the user doesn't have to know about, occured.
func UnexpectedErrorProblem(relPath string) *api.ProblemDetails {
	return &api.ProblemDetails{
		StatusCode: http.StatusInternalServerError,
		Title:      "Unexpected error occured",
		Detail:     "An unexpected error occured during processing your request",
		Instance:   relPath,
	}
}

// InvalidCredentials returns a problem details response used when a user enters
// invalid login credentials.
func InvalidCredentials(relPath string) *api.ProblemDetails {
	return &api.ProblemDetails{
		StatusCode: http.StatusUnauthorized,
		Title:      "Invalid credentials",
		Detail:     "Submitted credentials are invalid",
		Instance:   relPath,
	}
}

// GetProblemDetails is used when a error needs to be identified and user needs
// a specific problem details response corresponding to the identified error.
func GetProblemDetails(err error, relPath string) (problem *api.ProblemDetails) {

	problem = UnexpectedErrorProblem(relPath)
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		problem = mySQLProblem(*mysqlErr, relPath)
	} else if errors.Is(err, sql.ErrNoRows) {
		problem = sqlNoRows(relPath)
	}

	return problem
}

// BadRequest is used when user sends a invalid/malformed JSON request. Details
// are extracted from the error param, if the error param can't be casted as
// malformedRequest, generic UnexpectedErrorProblem is returned.
func BadRequest(err error, relPath string) *api.ProblemDetails {
	if malformedErr, ok := err.(malformedRequestErr); ok {
		return &api.ProblemDetails{
			StatusCode: malformedErr.status,
			Title:      "Bad request",
			Detail:     malformedErr.msg,
			Instance:   relPath,
		}
	}

	return UnexpectedErrorProblem(relPath)
}

// mySQLProblem function returns specific problem details accourding to the
// error number in the MySQLError. If the error has unknown number, generic
// UnexpectedErrorProblem is returned.
func mySQLProblem(err mysql.MySQLError, relPath string) *api.ProblemDetails {

	var statusCode int
	var title, detail string

	switch err.Number {
	case 1062:
		statusCode = http.StatusConflict
		title, detail = sqlDuplicateEntry(err)

	default:
		return UnexpectedErrorProblem(relPath)
	}

	return &api.ProblemDetails{
		StatusCode: statusCode,
		Title:      title,
		Detail:     detail,
		Instance:   relPath,
	}
}

// sqlDuplicateEntry is a util function for identifying which submitted entry is
// a duplicate according to the err. Then it returns corresponding type, title and
// detail for creation of a problem details response.
func sqlDuplicateEntry(err mysql.MySQLError) (title, detail string) {
	re := regexp.MustCompile(`'(.{3,30}?)'`)
	entry := re.FindString(err.Message)

	if strings.Contains(err.Message, "users.username") {
		title = "Username already exists"
		detail = fmt.Sprintf("Username %s already exists", entry)

	} else if strings.Contains(err.Message, "users.email") {
		title = "Email already used"
		detail = fmt.Sprintf("Email %s is already used", entry)

	} else {
		title = "Unknown duplicate entry"
		detail = fmt.Sprintf("Unknown entry %s is already used", entry)

	}

	return
}

// sqlNoRows is a util function for creating a problem details response when a
// submitted username was not found in database.
func sqlNoRows(relPath string) *api.ProblemDetails {
	title := "Entered username was not found"
	detail := "Username you entered was not found"

	return &api.ProblemDetails{
		StatusCode: http.StatusUnauthorized,
		Title:      title,
		Detail:     detail,
		Instance:   relPath,
	}
}
