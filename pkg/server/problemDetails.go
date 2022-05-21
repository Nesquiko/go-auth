package server

import (
	"net/http"

	"github.com/Nesquiko/go-auth/pkg/api"
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
