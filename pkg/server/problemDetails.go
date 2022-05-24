package server

import (
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

func BadRequest(err malformedRequest) api.ProblemDetails {
	return api.ProblemDetails{
		StatusCode: err.status,
		Type:       "bad.request",
		Title:      "Bad request",
		Detail:     err.msg,
	}
}

func MySQLProblem(err error) api.ProblemDetails {
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		var statusCode int
		var problemType, title, detail string

		switch mysqlErr.Number {
		case 1062:
			statusCode = http.StatusConflict
			problemType, title, detail = sqlDuplicateEntry(*mysqlErr)
		}

		return api.ProblemDetails{
			StatusCode: statusCode,
			Type:       problemType,
			Title:      title,
			Detail:     detail,
		}
	}

	return UnexpectedErrorProblem
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
