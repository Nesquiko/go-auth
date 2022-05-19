package server

import (
	"encoding/json"
	"net/http"

	"github.com/Nesquiko/go-auth/pkg/api"
	"github.com/Nesquiko/go-auth/pkg/db"
)

type GoAuthServer struct{}

func (s GoAuthServer) Signup(w http.ResponseWriter, r *http.Request) {
	connection, err := db.Connect("root", "goAuthDB", "users", "127.0.0.1:3306")
	if err != nil {
		respondWithError(w, UNEXPECTED_ERROR)
		return
	}

	var req api.SignupRequest
	err = decodeJSONBody(w, r, &req)
	if err != nil {
		panic(err)
	}

	hashedPassword, err := encryptPassword(req.Password)
	if err != nil {
		panic(err)
	}

	newUser := &db.UserModel{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
	}

	connection.SaveUser(newUser)
}

func (s GoAuthServer) Login(w http.ResponseWriter, r *http.Request) {

}

func respondWithError(w http.ResponseWriter, problem api.ProblemDetails) {

	w.Header().Set(CONTENT_TYPE, APPLICATION_JSON)
	w.WriteHeader(problem.StatusCode)

	json.NewEncoder(w).Encode(problem)
}
