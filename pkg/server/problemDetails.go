package server

import (
	"net/http"

	"github.com/Nesquiko/go-auth/pkg/api"
)

var UNEXPECTED_ERROR = api.ProblemDetails{
	StatusCode: http.StatusInternalServerError,
	Type:       "unexpected.error",
	Title:      "Unexpected error occured",
	Detail:     "An unexpected error occured during processing your request",
}
