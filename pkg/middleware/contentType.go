package middleware

import (
	"net/http"

	"github.com/Nesquiko/go-auth/pkg/api"
	"github.com/Nesquiko/go-auth/pkg/consts"
)

// ContentTypeFilter is a middleware for filtering requests which do not have
// Content-Type header set to application/json.
func ContentTypeFilter(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if ct := r.Header.Get(consts.ContentType); ct != consts.ApplicationJSON {
			responseMsg := "Content-Type header is not application/json"

			pd := api.ProblemDetails{
				StatusCode: http.StatusUnsupportedMediaType,
				Type:       "bad.request",
				Title:      "Bad request",
				Detail:     responseMsg,
			}

			respondWithProblemDetails(w, pd)
			return
		}

		next.ServeHTTP(w, r)
	})
}
