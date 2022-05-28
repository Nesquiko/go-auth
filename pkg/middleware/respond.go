package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/Nesquiko/go-auth/pkg/api"
	"github.com/Nesquiko/go-auth/pkg/consts"
)

// respondWithProblemDetails takes a problem details response created when an error
// occured during processing of a user request. It is serialized into a JSON.
// Then a status code is set to the one retrieved from problem details and
// a response is sent.
func respondWithProblemDetails(w http.ResponseWriter, problem api.ProblemDetails) {
	w.Header().Set(consts.ContentType, consts.ApplicationJSON)
	w.WriteHeader(problem.StatusCode)

	json.NewEncoder(w).Encode(problem)
}
