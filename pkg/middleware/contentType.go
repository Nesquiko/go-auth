package middleware

import (
	"net/http"

	"github.com/Nesquiko/go-auth/pkg/api"
	"github.com/Nesquiko/go-auth/pkg/consts"
)

// ContentTypeFilter is a middleware for filtering requests which do not have
// Content-Type header set to application/json. Firstly checks if any bearer
// token is in Headers, if yes, then proceed, if no check for correct Content-Type
// header.
func ContentTypeFilter(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if bearer := r.Header.Get(consts.Authorization); bearer != "" {
			next.ServeHTTP(w, r)
			return

		} else if ct := r.Header.Get(consts.ContentType); ct != consts.ApplicationJSON {
			responseMsg := "Content-Type header is not application/json"

			pd := api.ProblemDetails{
				StatusCode: http.StatusUnsupportedMediaType,
				Title:      "Bad request",
				Detail:     responseMsg,
				Instance:   r.URL.Path,
			}

			respondWithProblemDetails(w, pd)
			return
		}

		next.ServeHTTP(w, r)
	})
}
