// Package server provides functions for handling different API endpoints of the
// Go-Auth application.
package server

import (
	"crypto/rand"
	"encoding/base32"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Nesquiko/go-auth/pkg/api"
	"github.com/Nesquiko/go-auth/pkg/consts"
	"github.com/Nesquiko/go-auth/pkg/db"
	"github.com/Nesquiko/go-auth/pkg/security"
	"github.com/dgryski/dgoogauth"
)

// GoAuthServer is an empty struct used as a representation of a handler for
// API endpoints.
type GoAuthServer struct{}

// Signup handles when a user sends a request to the /signup endpoint for signing
// up. After successfully decoding JSON request, new user entry is saved into the
// database.
// Specific endpoint details can be found in ./openapi folder in the
// OpenAPI specification.
func (s GoAuthServer) Signup(w http.ResponseWriter, r *http.Request) {

	var req api.SignupRequest
	err := validateJSONRequestBody(w, r, &req)
	if err != nil {
		respondWithError(w, BadRequest(err, r.URL.Path))
		return
	}

	hashedPassword, err := security.EncryptPassword(req.Password)
	if err != nil {
		respondWithError(w, UnexpectedErrorProblem(r.URL.Path))
		return
	}

	newUser := &db.UserModel{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
	}

	err = db.DBConn.SaveUser(newUser)
	if err != nil {
		respondWithError(w, GetProblemDetails(err, r.URL.Path))
		return
	}
}

// Login handles when a user sends a request to the /login endpoint for logging
// in. After successfully decoding JSON request, user credentials are compared
// with corresponding ones retrieved from database. If credentials are valid
// new JWT token is generated for the user and sent.
// Specific endpoint details can be found in ./openapi folder in the
// OpenAPI specification.
func (s GoAuthServer) Login(w http.ResponseWriter, r *http.Request) {

	var req api.LoginRequest
	err := validateJSONRequestBody(w, r, &req)
	if err != nil {
		respondWithError(w, BadRequest(err, r.URL.Path))
		return
	}

	user, err := db.DBConn.UserByUsername(req.Username)
	if err != nil {
		respondWithError(w, GetProblemDetails(err, r.URL.Path))
		return
	}

	if !security.HashAndPasswordMatch(user.PasswordHash, req.Password) {
		respondWithError(w, InvalidCredentials(r.URL.Path))
		return
	}

	jwt, err := security.GenerateJWT(req.Username, false)
	if err != nil {
		respondWithError(w, UnexpectedErrorProblem(r.URL.Path))
		return
	}

	response := api.LoginResponse{UnauthToken: jwt}
	respondWithSuccess(w, response)
}

// Setup2FA creates new 2FA secret for user and returns a 2FA uri
// for generating QR code.
func (s GoAuthServer) Setup2FA(w http.ResponseWriter, r *http.Request) {

	bearer := r.Header.Get(consts.Authorization)
	if bearer == "" {
		respondWithError(w, Unauthorized(r.URL.Path))
		return
	}

	token := strings.Split(bearer, consts.BearerPrefix)[1]
	c, err := security.ValidateToken(token)
	if err != nil {
		respondWithError(w, Unauthorized(r.URL.Path))
		return
	}

	enabled, err := db.DBConn.GetEnabled2FA(c.Username)
	if err != nil {
		respondWithError(w, UnexpectedErrorProblem(r.URL.Path))
		return
	}

	if enabled {
		respondWithError(w, Unauthorized(r.URL.Path))
		return
	}

	random := make([]byte, 10)
	rand.Read(random)
	secret := base32.StdEncoding.EncodeToString(random)

	err = db.DBConn.Save2FASecret(c.Username, secret)
	if err != nil {
		respondWithError(w, UnexpectedErrorProblem(r.URL.Path))
		return
	}
	authLink := fmt.Sprintf(
		"otpauth://totp/GoAuth:%s?secret=%s&issuer=GoAuth",
		c.Username,
		secret,
	)
	response := api.Secret2FAResponse{QrURI: &authLink}
	respondWithSuccess(w, response)
}

func (s GoAuthServer) Verify2FA(w http.ResponseWriter, r *http.Request) {

	bearer := r.Header.Get(consts.Authorization)
	if bearer == "" {
		respondWithError(w, Unauthorized(r.URL.Path))
		return
	}

	token := strings.Split(bearer, consts.BearerPrefix)[1]
	c, err := security.ValidateToken(token)
	if err != nil {
		respondWithError(w, Unauthorized(r.URL.Path))
		return
	}

	sec, _ := db.DBConn.Get2FASecret(c.Username)
	otpc := &dgoogauth.OTPConfig{
		Secret:      sec,
		WindowSize:  3,
		HotpCounter: 0,
	}
	var req api.Verify2FAJSONRequestBody
	err = validateJSONRequestBody(w, r, &req)
	if err != nil {
		respondWithError(w, BadRequest(err, r.URL.Path))
		return
	}

	ok, err := otpc.Authenticate(strconv.FormatInt(int64(req.Otp), 10))
	if err != nil {
		respondWithError(w, UnexpectedErrorProblem(r.URL.Path))
		return
	}

	if !ok {
		respondWithError(w, Unauthorized(r.URL.Path))
		return
	}

	jwt, err := security.GenerateJWT(c.Username, true)
	if err != nil {
		respondWithError(w, UnexpectedErrorProblem(r.URL.Path))
		return
	}
	response := api.VerifyResponse{
		AccessToken: jwt,
	}

	enabled, err := db.DBConn.GetEnabled2FA(c.Username)
	if err != nil {
		respondWithError(w, UnexpectedErrorProblem(r.URL.Path))
		return
	}

	if !enabled {
		db.DBConn.UpdateEnabled2FA(c.Username, true)
	}

	respondWithSuccess(w, response)
}

func (s GoAuthServer) TestAuth(w http.ResponseWriter, r *http.Request) {
	bearer := r.Header.Get(consts.Authorization)
	if bearer == "" {
		respondWithError(w, Unauthorized(r.URL.Path))
		return
	}

	token := strings.Split(bearer, consts.BearerPrefix)[1]
	c, err := security.ValidateToken(token)
	if err != nil {
		respondWithError(w, Unauthorized(r.URL.Path))
		return
	}

	if !c.Authenticated {
		respondWithError(w, Unauthorized(r.URL.Path))
		return
	}

	respondWithSuccess(w, "Full access granted")
}

// respondWithSuccess takes a response to be returned to a user making a
// request and serializes it into a JSON. Then sets a http.StatusOK as the
// response status code and then the response is sent to user.
func respondWithSuccess[T any](w http.ResponseWriter, response T) {
	w.Header().Set(consts.ContentType, consts.ApplicationJSON)
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(response)
}

// respondWithError takes a problem details response created when an error
// occured during processing of a user request. It is serialized into a JSON.
// Then a status code is set to the one retrieved from problem details and
// a response is sent
func respondWithError(w http.ResponseWriter, problem *api.ProblemDetails) {
	w.Header().Set(consts.ContentType, consts.ApplicationJSON)
	w.WriteHeader(problem.StatusCode)

	json.NewEncoder(w).Encode(problem)
}
