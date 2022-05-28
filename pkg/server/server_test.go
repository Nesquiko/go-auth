package server

import (
	"bytes"
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
	"github.com/go-chi/chi/v5"
)

var server http.Handler

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
		name                            string
		reqString                       string
		wantCode                        int
		wantType, wantTitle, wantDetail string
	}{
		{
			"BadlyFormedJSONBodyAtPosition",
			// missing , right 	     here
			"{\"email\":\"test@foo.com\"\"password\":\"foobarz\",\"username\":\"Barz\"}",
			http.StatusBadRequest, "bad.request", "Bad request",
			fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", 24),
		},
		{
			"BadlyFormedJSONBody",
			// missing } at the end
			"{\"email\":\"test@foo.com\",\"password\":\"foobarz\",\"username\":\"Barz\"",
			http.StatusBadRequest, "bad.request", "Bad request",
			"Request body contains badly-formed JSON",
		},
		{
			"InvalidValueForField",
			"{\"email\":123,\"password\":\"foobarz\",\"username\":\"Barz\"}",
			http.StatusBadRequest, "bad.request", "Bad request",
			fmt.Sprintf(
				"Request body contains an invalid value for the %q field (at position %d)",
				"email",
				12,
			),
		},
		{
			"UnknownField",
			"{\"unknown\":\"field\",\"email\":\"test@foo.com\",\"password\":\"foobarz\",\"username\":\"Barz\"}",
			http.StatusBadRequest, "bad.request", "Bad request",
			fmt.Sprintf("Request body contains unknown field %q", "unknown"),
		},
		{
			"EmptyRequest",
			"",
			http.StatusBadRequest, "bad.request", "Bad request",
			"Request body must not be empty",
		},
		{
			"LargeBody",
			fmt.Sprintf("{\"email\":%q,\"password\":\"foobarz\",\"username\":\"Barz\"}",
				strings.Repeat("email", 140)),
			http.StatusRequestEntityTooLarge, "bad.request", "Bad request",
			fmt.Sprintf("Request body must not be larger than %dB", maxSize),
		},
		{
			"MissingField",
			"{\"email\":\"test@foo.com\",\"password\":\"foobarz\"}",
			http.StatusBadRequest, "bad.request", "Bad request",
			"Request body is not complete",
		},
	}

	for _, tc := range testCases {
		req := httptest.NewRequest("POST", "/signup", strings.NewReader(tc.reqString))
		req.Header.Add(consts.ContentType, consts.ApplicationJSON)

		wantBody := fmt.Sprintf("{%q:%d,%q:%q,%q:%q,%q:%q}\n",
			"status_code", tc.wantCode,
			"type", tc.wantType,
			"title", tc.wantTitle,
			"detail", tc.wantDetail)

		res := executeRequest(req)

		if res.Code != tc.wantCode {
			t.Errorf("Expected status code to be %d, but was %d", tc.wantCode, res.Code)
		}
		if res.Body.String() != wantBody {
			t.Errorf("Expected response body to be %s, but was %s", wantBody, res.Body)
		}
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
	wantType := "username.already_exists"
	wantTitle := "Username already exists"
	wantDetail := fmt.Sprintf("Username '%s' already exists", username)

	wantBody := fmt.Sprintf("{%q:%d,%q:%q,%q:%q,%q:%q}\n",
		"status_code", wantCode,
		"type", wantType,
		"title", wantTitle,
		"detail", wantDetail)

	res := executeRequest(req)

	if res.Code != wantCode {
		t.Errorf("Expected status code to be %d, but was %d", wantCode, res.Code)
	}
	if res.Body.String() != wantBody {
		t.Errorf("Expected response body to be %s, but was %s", wantBody, res.Body)
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
	wantType := "email.already_used"
	wantTitle := "Email already used"
	wantDetail := fmt.Sprintf("Email '%s' is already used", email)

	wantBody := fmt.Sprintf("{%q:%d,%q:%q,%q:%q,%q:%q}\n",
		"status_code", wantCode,
		"type", wantType,
		"title", wantTitle,
		"detail", wantDetail)

	res := executeRequest(req)

	if res.Code != wantCode {
		t.Errorf("Expected status code to be %d, but was %d", wantCode, res.Code)
	}
	if res.Body.String() != wantBody {
		t.Errorf("Expected response body to be %s, but was %s", wantBody, res.Body)
	}
}
