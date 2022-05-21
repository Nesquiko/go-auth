package server

import (
	"encoding/json"
	"net/http"

	"github.com/Nesquiko/go-auth/pkg/api"
	"github.com/Nesquiko/go-auth/pkg/db"
)

type GoAuthServer struct{}

func (s GoAuthServer) Signup(w http.ResponseWriter, r *http.Request) {
	dbCon := db.DBConnection

	var req api.SignupRequest
	err := decodeJSONBody(w, r, &req)
	if err != nil {
		respondWithError(w, BadRequest(err.(malformedRequest)))
		return
	}

	hashedPassword, err := encryptPassword(req.Password)
	if err != nil {
		respondWithError(w, UnexpectedErrorProblem)
		return
	}

	newUser := &db.UserModel{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
	}

	err = dbCon.SaveUser(newUser)
	if err != nil {
		respondWithError(w, SQLProblem(err))
		return
	}
}

func (s GoAuthServer) Login(w http.ResponseWriter, r *http.Request) {

}

func respondWithError(w http.ResponseWriter, problem api.ProblemDetails) {
	w.Header().Set(contentType, applicationJSON)
	w.WriteHeader(problem.StatusCode)

	json.NewEncoder(w).Encode(problem)
}
