package server

import (
	"encoding/json"
	"net/http"

	"github.com/Nesquiko/go-auth/pkg/api"
	"github.com/Nesquiko/go-auth/pkg/db"
	"github.com/Nesquiko/go-auth/pkg/security"
)

type GoAuthServer struct{}

func (s GoAuthServer) Signup(w http.ResponseWriter, r *http.Request) {
	dbConn := db.DBConnection

	var req api.SignupRequest
	err := decodeJSONBody(w, r, &req)
	if err != nil {
		respondWithError(w, BadRequest(err))
		return
	}

	hashedPassword, err := security.EncryptPassword(req.Password)
	if err != nil {
		respondWithError(w, UnexpectedErrorProblem)
		return
	}

	newUser := &db.UserModel{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
	}

	err = dbConn.SaveUser(newUser)
	if err != nil {
		respondWithError(w, GetProblemDetails(err))
		return
	}
}

func (s GoAuthServer) Login(w http.ResponseWriter, r *http.Request) {
	dbCon := db.DBConnection

	var req api.LoginRequest
	err := decodeJSONBody(w, r, &req)
	if err != nil {
		respondWithError(w, UnexpectedErrorProblem)
		return
	}

	user, err := dbCon.UserByUsername(req.Username)
	if err != nil {
		respondWithError(w, GetProblemDetails(err))
		return
	}

	if !security.HashAndPasswordMatch(user.PasswordHash, req.Password) {
		respondWithError(w, InvalidCredentials)
		return
	}

	jwt, err := security.GenerateJWT(req.Username)
	if err != nil {
		respondWithError(w, UnexpectedErrorProblem)
		return
	}

	response := api.LoginResponse{AccessToken: jwt}
	respondWithSuccess(w, response)
}

func respondWithSuccess[T any](w http.ResponseWriter, response T) {
	w.Header().Set(contentType, applicationJSON)
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(response)
}

func respondWithError(w http.ResponseWriter, problem api.ProblemDetails) {
	w.Header().Set(contentType, applicationJSON)
	w.WriteHeader(problem.StatusCode)

	json.NewEncoder(w).Encode(problem)
}
