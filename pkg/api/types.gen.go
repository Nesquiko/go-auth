// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.11.1-0.20220526233757-b1edf760c6e7 DO NOT EDIT.
package api

const (
	UnauthBearerTokenScopes = "unauthBearerToken.Scopes"
)

// A problem details response, which occured during processing of a request. (Trying to adhere to RFC 7807)
type ProblemDetails struct {
	// Human-readable explanation specific to this occurrence of the problem
	Detail string `json:"detail"`

	// A URI reference that identifies the specific occurrence of the problem
	Instance string `json:"instance"`

	// A http status code describing a problem
	StatusCode int `json:"status_code"`

	// A short, human-readable summary of the problem type
	Title string `json:"title"`
}

// LoginResponse defines model for LoginResponse.
type LoginResponse struct {
	// An unauthenticated JWT access token needed in 2FA.
	UnauthToken string `json:"unauth_token"`
}

// Secret2FAResponse defines model for Secret2FAResponse.
type Secret2FAResponse struct {
	QrURI *string `json:"qrURI,omitempty"`
}

// VerifyResponse defines model for VerifyResponse.
type VerifyResponse struct {
	// An full access JWT.
	AccessToken string `json:"access_token"`
}

// LoginRequest defines model for LoginRequest.
type LoginRequest struct {
	// Password of an user account
	Password string `json:"password" validate:"required"`

	// Username of an user account
	Username string `json:"username" validate:"required"`
}

// SignupRequest defines model for SignupRequest.
type SignupRequest struct {
	// Email address of a new user account
	Email string `json:"email" validate:"required"`

	// Password for getting access to the new user account
	Password string `json:"password" validate:"required"`

	// Username with which new user account will be identified in the system
	Username string `json:"username" validate:"required"`
}

// LoginJSONRequestBody defines body for Login for application/json ContentType.
type LoginJSONRequestBody LoginRequest

// SignupJSONRequestBody defines body for Signup for application/json ContentType.
type SignupJSONRequestBody SignupRequest
