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

var UnexpectedErrorProblem = api.ProblemDetails{
	StatusCode: http.StatusInternalServerError,
	Type:       "unexpected.error",
	Title:      "Unexpected error occured",
	Detail:     "An unexpected error occured during processing your request",
}

var InvalidCredentials = api.ProblemDetails{
	StatusCode: http.StatusUnauthorized,
	Type:       "credentials.invalid",
	Title:      "Invalid credentials",
	Detail:     "Submitted credentials are invalid",
}

func GetProblemDetails(err error) (problem api.ProblemDetails) {

	problem = UnexpectedErrorProblem
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		problem = mySQLProblem(*mysqlErr)
	} else if errors.Is(err, sql.ErrNoRows) {
		problem = sqlNoRows()
	}

	return problem
}

func BadRequest(err error) api.ProblemDetails {
	if malformedErr, ok := err.(malformedRequest); ok {
		return api.ProblemDetails{
			StatusCode: malformedErr.status,
			Type:       "bad.request",
			Title:      "Bad request",
			Detail:     malformedErr.msg,
		}
	}

	return UnexpectedErrorProblem
}

func mySQLProblem(err mysql.MySQLError) api.ProblemDetails {

	var statusCode int
	var problemType, title, detail string

	switch err.Number {
	case 1062:
		statusCode = http.StatusConflict
		problemType, title, detail = sqlDuplicateEntry(err)

	default:
		return UnexpectedErrorProblem
	}

	return api.ProblemDetails{
		StatusCode: statusCode,
		Type:       problemType,
		Title:      title,
		Detail:     detail,
	}
}

func sqlDuplicateEntry(err mysql.MySQLError) (problemType, title, detail string) {
	re := regexp.MustCompile(`'(.{3,30}?)'`)
	entry := re.FindString(err.Message)

	if strings.Contains(err.Message, "users.username") {
		problemType = "username.already_exists"
		title = "Username already exists"
		detail = fmt.Sprintf("Username %s already exists", entry)

	} else if strings.Contains(err.Message, "users.email") {
		problemType = "email.already_used"
		title = "Email already used"
		detail = fmt.Sprintf("Email %s is already used", entry)

	} else {
		problemType = "unknown.duplicate"
		title = "Unknown duplicate entry"
		detail = fmt.Sprintf("Unknown entry %s is already used", entry)

	}

	return
}

func sqlNoRows() api.ProblemDetails {
	problemType := "username.not_found"
	title := "Entered username was not found"
	detail := "Username you eneterd was not found"

	return api.ProblemDetails{
		StatusCode: http.StatusUnauthorized,
		Type:       problemType,
		Title:      title,
		Detail:     detail,
	}
}
