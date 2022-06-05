package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
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

			wantBody := fmt.Sprintf("{%q:%d,%q:%q,%q:%q,%q:%q}\n",
				"status_code", wantCode,
				"title", wantTitle,
				"detail", wantDetail,
				"instance", wantInstance,
			)

			res := executeRequest(req)

			if res.Code != wantCode {
				t.Errorf("Expected status code to be %d, but was %d", wantCode, res.Code)
			}
			if res.Body.String() != wantBody {
				t.Errorf("Expected response body to be %s, but was %s", wantBody, res.Body)
			}
		})
	}
}
