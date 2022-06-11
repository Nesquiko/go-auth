package middleware

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Nesquiko/go-auth/pkg/api"
	"github.com/Nesquiko/go-auth/pkg/consts"
	"github.com/Nesquiko/go-auth/pkg/server"
	"github.com/go-chi/chi/v5"
)

var handler http.Handler

func TestMain(m *testing.M) {
	r := chi.NewRouter()
	var s server.GoAuthServer

	middlewares := []api.MiddlewareFunc{
		ContentTypeFilter,
	}
	servOpts := api.ChiServerOptions{
		BaseRouter:  r,
		Middlewares: middlewares,
	}

	handler = api.HandlerWithOptions(s, servOpts)

	code := m.Run()

	os.Exit(code)
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	return rr
}

func TestContentTypeFilterInvalidHeader(t *testing.T) {
	testCases := []struct {
		name          string
		reqString     string
		contentHeader string
	}{
		{"NoContentTypeHeader", "\"{}\"", ""},
		{"WrongContentTypeHeader", "\"{}\"", "textPlain"},
	}

	for i, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := json.NewEncoder(&buf).Encode(tc.reqString)
			if err != nil {
				t.Fatalf("Error in encoding of request in test case %d", i)
			}

			req := httptest.NewRequest("POST", "/signup", &buf)
			req.Header.Add(consts.ContentType, tc.contentHeader)

			wantCode := http.StatusUnsupportedMediaType
			wantTitle := "Bad request"
			wantDetail := "Content-Type header is not application/json"
			wantInstance := "/signup"

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
		})
	}
}
