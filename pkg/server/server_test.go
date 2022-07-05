package server

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/Nesquiko/go-auth/pkg/api"
	"github.com/Nesquiko/go-auth/pkg/consts"
	"github.com/Nesquiko/go-auth/pkg/db"
	"github.com/Nesquiko/go-auth/pkg/db/mocks"
	"github.com/Nesquiko/go-auth/pkg/security"
	"github.com/go-chi/chi/v5"
)

var server http.Handler
var signupPath = "/signup"
var loginPath = "/login"

func TestMain(m *testing.M) {
	r := chi.NewRouter()
	var s GoAuthServer
	servOpts := api.ChiServerOptions{
		BaseRouter: r,
	}
	server = api.HandlerWithOptions(s, servOpts)

	db.DBConn = mocks.DBConnectionMock{}

	code := m.Run()

	os.Exit(code)
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	server.ServeHTTP(rr, req)

	return rr
}

func TestSignupBadRequest(t *testing.T) {
	testCases := []struct {
		name                                string
		reqString                           string
		wantCode                            int
		wantTitle, wantDetail, wantInstance string
	}{
		{
			"BadlyFormedJSONBodyAtPosition",
			// missing , right 	     here
			"{\"email\":\"test@foo.com\"\"password\":\"foobarz\",\"username\":\"Barz\"}",
			http.StatusBadRequest, "Bad request",
			fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", 24),
			signupPath,
		},
		{
			"BadlyFormedJSONBody",
			// missing } at the end
			"{\"email\":\"test@foo.com\",\"password\":\"foobarz\",\"username\":\"Barz\"",
			http.StatusBadRequest, "Bad request",
			"Request body contains badly-formed JSON",
			signupPath,
		},
		{
			"InvalidValueForField",
			"{\"email\":123,\"password\":\"foobarz\",\"username\":\"Barz\"}",
			http.StatusBadRequest, "Bad request",
			fmt.Sprintf(
				"Request body contains an invalid value for the %q field (at position %d)",
				"email",
				12,
			),
			signupPath,
		},
		{
			"UnknownField",
			"{\"unknown\":\"field\",\"email\":\"test@foo.com\",\"password\":\"foobarz\",\"username\":\"Barz\"}",
			http.StatusBadRequest, "Bad request",
			fmt.Sprintf("Request body contains unknown field %q", "unknown"),
			signupPath,
		},
		{
			"EmptyRequest",
			"",
			http.StatusBadRequest, "Bad request",
			"Request body must not be empty",
			signupPath,
		},
		{
			"LargeBody",
			fmt.Sprintf("{\"email\":%q,\"password\":\"foobarz\",\"username\":\"Barz\"}",
				strings.Repeat("email", 140)),
			http.StatusRequestEntityTooLarge, "Bad request",
			fmt.Sprintf("Request body must not be larger than %dB", maxSize),
			signupPath,
		},
		{
			"MissingField",
			"{\"email\":\"test@foo.com\",\"password\":\"foobarz\"}",
			http.StatusBadRequest, "Bad request",
			"Request body is not complete",
			signupPath,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", signupPath, strings.NewReader(tc.reqString))
			req.Header.Add(consts.ContentType, consts.ApplicationJSON)

			res := executeRequest(req)
			var pd api.ProblemDetails
			err := json.Unmarshal(res.Body.Bytes(), &pd)
			if err != nil {
				t.Fatalf("Error when unmarshalling: %s", err.Error())
			}

			if pd.StatusCode != tc.wantCode {
				t.Errorf("Status code, expected %d, but was %d", tc.wantCode, pd.StatusCode)
			}
			if pd.Title != tc.wantTitle {
				t.Errorf("Title, expected %q, but was %q", tc.wantTitle, pd.Title)
			}
			if pd.Detail != tc.wantDetail {
				t.Errorf("Detail, expected %q, but was %q", tc.wantDetail, pd.Detail)
			}
			if pd.Instance != tc.wantInstance {
				t.Errorf("Instance, expected %q, but was %q", tc.wantInstance, pd.Instance)
			}
		})
	}
}

func TestSignupValidRequest(t *testing.T) {
	reqBody := api.SignupRequest{
		Email:    "test@foo.com",
		Username: "Barz",
		Password: "foobarz",
	}
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(reqBody)

	if err != nil {
		t.Fatal("Error in encoding of struct")
	}

	req := httptest.NewRequest("POST", "/signup", &buf)
	req.Header.Add(consts.ContentType, consts.ApplicationJSON)

	wantCode := http.StatusOK

	res := executeRequest(req)

	if res.Code != wantCode {
		t.Errorf("Expected status code to be %d, but was %d", wantCode, res.Code)
	}
}

func TestSignupUsernameAlreadyExists(t *testing.T) {
	username := "Barz"
	db.DBConn.SaveUser(&db.UserModel{
		Email:        "bar@foo.com",
		Username:     username,
		PasswordHash: "hash",
	})

	reqBody := api.SignupRequest{
		Email:    "test@foo.com",
		Username: username,
		Password: "foobarz",
	}
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(reqBody)

	if err != nil {
		t.Fatal("Error in encoding of struct")
	}

	req := httptest.NewRequest("POST", "/signup", &buf)
	req.Header.Add(consts.ContentType, consts.ApplicationJSON)

	wantCode := http.StatusConflict
	wantTitle := "Username already exists"
	wantDetail := fmt.Sprintf("Username '%s' already exists", username)
	wantInstance := signupPath

	res := executeRequest(req)
	var pd api.ProblemDetails
	err = json.Unmarshal(res.Body.Bytes(), &pd)
	if err != nil {
		t.Fatalf("Error when unmarshalling: %s", err.Error())
	}

	if pd.StatusCode != wantCode {
		t.Errorf("Status code, expected %d, but was %d", wantCode, pd.StatusCode)
	}
	if pd.Title != wantTitle {
		t.Errorf("Title, expected %q, but was %q", wantTitle, pd.Title)
	}
	if pd.Detail != wantDetail {
		t.Errorf("Detail, expected %q, but was %q", wantDetail, pd.Detail)
	}
	if pd.Instance != wantInstance {
		t.Errorf("Instance, expected %q, but was %q", wantInstance, pd.Instance)
	}
}

func TestSignupEmailAlreadyExists(t *testing.T) {
	email := "bar@foo.com"
	db.DBConn.SaveUser(&db.UserModel{
		Email:        email,
		Username:     "Bar",
		PasswordHash: "hash",
	})

	reqBody := api.SignupRequest{
		Email:    email,
		Username: "John",
		Password: "foobarz",
	}
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(reqBody)

	if err != nil {
		t.Fatal("Error in encoding of struct")
	}

	req := httptest.NewRequest("POST", "/signup", &buf)
	req.Header.Add(consts.ContentType, consts.ApplicationJSON)

	wantCode := http.StatusConflict
	wantTitle := "Email already used"
	wantDetail := fmt.Sprintf("Email '%s' is already used", email)
	wantInstance := signupPath

	res := executeRequest(req)
	var pd api.ProblemDetails
	err = json.Unmarshal(res.Body.Bytes(), &pd)
	if err != nil {
		t.Fatalf("Error when unmarshalling: %s", err.Error())
	}

	if pd.StatusCode != wantCode {
		t.Errorf("Status code, expected %d, but was %d", wantCode, pd.StatusCode)
	}
	if pd.Title != wantTitle {
		t.Errorf("Title, expected %q, but was %q", wantTitle, pd.Title)
	}
	if pd.Detail != wantDetail {
		t.Errorf("Detail, expected %q, but was %q", wantDetail, pd.Detail)
	}
	if pd.Instance != wantInstance {
		t.Errorf("Instance, expected %q, but was %q", wantInstance, pd.Instance)
	}
}

func TestLoginBadRequest(t *testing.T) {
	testCases := []struct {
		name                                string
		reqString                           string
		wantCode                            int
		wantTitle, wantDetail, wantInstance string
	}{
		{
			"BadlyFormedJSONBodyAtPosition",
			// missing , right 	     here
			"{\"username\":\"Barz\"\"password\":\"foobarz\"}",
			http.StatusBadRequest, "Bad request",
			fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", 19),
			loginPath,
		},
		{
			"BadlyFormedJSONBody",
			// missing } at the end
			"{\"username\":\"Barz\",\"password\":\"foobarz\"",
			http.StatusBadRequest, "Bad request",
			"Request body contains badly-formed JSON",
			loginPath,
		},
		{
			"InvalidValueForField",
			"{\"username\":123,\"password\":\"foobarz\"}",
			http.StatusBadRequest, "Bad request",
			fmt.Sprintf(
				"Request body contains an invalid value for the %q field (at position %d)",
				"username",
				15,
			),
			loginPath,
		},
		{
			"UnknownField",
			"{\"username\":\"Barz\",\"password\":\"foobarz\",\"unknown\":\"field\"}",
			http.StatusBadRequest, "Bad request",
			fmt.Sprintf("Request body contains unknown field %q", "unknown"),
			loginPath,
		},
		{
			"EmptyRequest",
			"",
			http.StatusBadRequest, "Bad request",
			"Request body must not be empty",
			loginPath,
		},
		{
			"LargeBody",
			fmt.Sprintf("{\"username\":%q,\"password\":\"foobarz\"}",
				strings.Repeat("Josh", 150)),
			http.StatusRequestEntityTooLarge, "Bad request",
			fmt.Sprintf("Request body must not be larger than %dB", maxSize),
			loginPath,
		},
		{
			"MissingField",
			"{\"username\":\"Barz\"}",
			http.StatusBadRequest, "Bad request",
			"Request body is not complete",
			loginPath,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/login", strings.NewReader(tc.reqString))
			req.Header.Add(consts.ContentType, consts.ApplicationJSON)

			res := executeRequest(req)
			var pd api.ProblemDetails
			err := json.Unmarshal(res.Body.Bytes(), &pd)
			if err != nil {
				t.Fatalf("Error when unmarshalling: %s", err.Error())
			}

			if pd.StatusCode != tc.wantCode {
				t.Errorf("Status code, expected %d, but was %d", tc.wantCode, pd.StatusCode)
			}
			if pd.Title != tc.wantTitle {
				t.Errorf("Title, expected %q, but was %q", tc.wantTitle, pd.Title)
			}
			if pd.Detail != tc.wantDetail {
				t.Errorf("Detail, expected %q, but was %q", tc.wantDetail, pd.Detail)
			}
			if pd.Instance != tc.wantInstance {
				t.Errorf("Instance, expected %q, but was %q", tc.wantInstance, pd.Instance)
			}
		})
	}
}

func TestLoginUsernameDoesntExists(t *testing.T) {
	username := "James"

	reqBody := api.LoginRequest{
		Password: "password",
		Username: username,
	}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(reqBody)

	if err != nil {
		t.Fatal("Error in encoding of struct")
	}

	req := httptest.NewRequest("POST", "/login", &buf)
	req.Header.Add(consts.ContentType, consts.ApplicationJSON)

	wantCode := http.StatusUnauthorized
	wantTitle := "Invalid credentials"
	wantDetail := "Submitted credentials are invalid"
	wantInstance := loginPath

	res := executeRequest(req)
	var pd api.ProblemDetails
	err = json.Unmarshal(res.Body.Bytes(), &pd)
	if err != nil {
		t.Fatalf("Error when unmarshalling: %s", err.Error())
	}

	if pd.StatusCode != wantCode {
		t.Errorf("Status code, expected %d, but was %d", wantCode, pd.StatusCode)
	}
	if pd.Title != wantTitle {
		t.Errorf("Title, expected %q, but was %q", wantTitle, pd.Title)
	}
	if pd.Detail != wantDetail {
		t.Errorf("Detail, expected %q, but was %q", wantDetail, pd.Detail)
	}
	if pd.Instance != wantInstance {
		t.Errorf("Instance, expected %q, but was %q", wantInstance, pd.Instance)
	}
}

func TestLoginInvalidPassword(t *testing.T) {
	username := "James"
	passwd := "123"
	passwordHash, _ := security.EncryptPassword(passwd)

	db.DBConn.SaveUser(&db.UserModel{
		Email:        "james@barz.com",
		Username:     username,
		PasswordHash: passwordHash,
	})

	reqBody := api.LoginRequest{
		Password: passwd + "4",
		Username: username,
	}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(reqBody)

	if err != nil {
		t.Fatal("Error in encoding of struct")
	}

	req := httptest.NewRequest("POST", "/login", &buf)
	req.Header.Add(consts.ContentType, consts.ApplicationJSON)

	wantCode := http.StatusUnauthorized
	wantTitle := "Invalid credentials"
	wantDetail := "Submitted credentials are invalid"
	wantInstance := loginPath

	res := executeRequest(req)
	var pd api.ProblemDetails
	err = json.Unmarshal(res.Body.Bytes(), &pd)
	if err != nil {
		t.Fatalf("Error when unmarshalling: %s", err.Error())
	}

	if pd.StatusCode != wantCode {
		t.Errorf("Status code, expected %d, but was %d", wantCode, pd.StatusCode)
	}
	if pd.Title != wantTitle {
		t.Errorf("Title, expected %q, but was %q", wantTitle, pd.Title)
	}
	if pd.Detail != wantDetail {
		t.Errorf("Detail, expected %q, but was %q", wantDetail, pd.Detail)
	}
	if pd.Instance != wantInstance {
		t.Errorf("Instance, expected %q, but was %q", wantInstance, pd.Instance)
	}
}

func TestLoginValidRequest(t *testing.T) {
	username := "Susan"
	passwd := "123"
	passwordHash, _ := security.EncryptPassword(passwd)

	db.DBConn.SaveUser(&db.UserModel{
		Email:        "susan@barz.com",
		Username:     username,
		PasswordHash: passwordHash,
	})

	reqBody := api.LoginRequest{
		Password: passwd,
		Username: username,
	}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(reqBody)

	if err != nil {
		t.Fatal("Error in encoding of struct")
	}

	req := httptest.NewRequest("POST", "/login", &buf)
	req.Header.Add(consts.ContentType, consts.ApplicationJSON)

	wantCode := http.StatusOK
	wantClaimUsername := fmt.Sprintf("%q:%q", "username", username)
	var resBody api.LoginResponse

	res := executeRequest(req)
	json.Unmarshal(res.Body.Bytes(), &resBody)

	if res.Code != wantCode {
		t.Errorf("Expected status code to be %d, but was %d", wantCode, res.Code)
	}

	jwtSplit := strings.Split(resBody.UnauthToken, ".")
	payload, err := base64.RawURLEncoding.DecodeString(jwtSplit[1])
	if err != nil {
		t.Fatalf("decoding err was not nil, %q", err.Error())
	}
	if !strings.Contains(string(payload), wantClaimUsername) {
		t.Errorf("No valid username claim in payload %s", payload)
	}

}
