package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/Nesquiko/go-auth/pkg/consts"
)

type testJSONStruct struct {
	FieldString string `json:"FieldString" validate:"required"`
	FieldInt    int    `json:"FieldInt"    validate:"required"`
}

func Test_decodeJSONBodyValidJSONBody(t *testing.T) {
	var js testJSONStruct
	fieldString, fieldInt := "value", 123
	reqBody := testJSONStruct{fieldString, fieldInt}
	reqBodyJSON, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("", "", bytes.NewReader(reqBodyJSON))
	req.Header.Add(consts.ContentType, consts.ApplicationJSON)
	w := httptest.NewRecorder()

	err := validateJSONRequestBody(w, req, &js)

	if err != nil {
		t.Fatalf("Error that occured: %s", err.Error())
	}
	if js.FieldString != fieldString {
		t.Errorf("Expected string field to be %d, but was %d", fieldInt, js.FieldInt)
	}
	if js.FieldInt != fieldInt {
		t.Errorf("Expected int field to be %d, but was %d", fieldInt, js.FieldInt)
	}
}

func Test_decodeJSONBodyBadlyFormedJSONBodyAtPosition(t *testing.T) {
	var js testJSONStruct
	fieldString, fieldInt := "value", 123
	reqBody := testJSONStruct{fieldString, fieldInt}
	reqBodyJSON, _ := json.Marshal(reqBody)

	// , at position 23 is removed
	badlyFormedJSON := strings.Replace(string(reqBodyJSON), ",", "", 1)

	wantCode := http.StatusBadRequest
	wantPos := 23
	wantMsg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", wantPos)

	req, _ := http.NewRequest("", "", strings.NewReader(badlyFormedJSON))
	req.Header.Add(consts.ContentType, consts.ApplicationJSON)
	w := httptest.NewRecorder()

	err := validateJSONRequestBody(w, req, &js)
	if err == nil {
		t.Fatal("Error was nil")
	}

	mErr, ok := err.(malformedRequestErr)
	if !ok {
		t.Fatal("Error is not malformedRequest, but it should")
	}

	re := regexp.MustCompile(`\d{1,5}`)
	pos, _ := strconv.Atoi(re.FindString(mErr.msg))
	if wantPos != pos {
		t.Errorf("Expected position %d, but was %d", wantPos, pos)
	}
	if mErr.status != wantCode {
		t.Errorf("Wrong http status code, expected %d, but was %d", wantCode, mErr.status)
	}
	if mErr.msg != wantMsg {
		t.Errorf("Expected %s\nbut was %s", wantMsg, mErr.msg)
	}
}

func Test_decodeJSONBodyBadlyFormedJSONBody(t *testing.T) {
	var js testJSONStruct
	fieldString, fieldInt := "value", 123
	reqBody := testJSONStruct{fieldString, fieldInt}
	reqBodyJSON, _ := json.Marshal(reqBody)

	// remove } at the end
	badlyFormedJSON := strings.Replace(string(reqBodyJSON), "}", "", 1)

	wantCode := http.StatusBadRequest
	wantMsg := "Request body contains badly-formed JSON"

	req, _ := http.NewRequest("", "", strings.NewReader(badlyFormedJSON))
	req.Header.Add(consts.ContentType, consts.ApplicationJSON)
	w := httptest.NewRecorder()

	err := validateJSONRequestBody(w, req, &js)
	if err == nil {
		t.Fatal("Error was nil")
	}

	mErr, ok := err.(malformedRequestErr)
	if !ok {
		t.Fatal("Error is not malformedRequest, but it should")
	}

	if mErr.status != wantCode {
		t.Errorf("Wrong http status code, expected %d, but was %d", wantCode, mErr.status)
	}
	if mErr.msg != wantMsg {
		t.Errorf("Expected %s\nbut was %s", wantMsg, mErr.msg)
	}
}

func Test_decodeJSONBodyInvalidValueForField(t *testing.T) {
	var js testJSONStruct
	fieldString, fieldInt := "value", 123
	reqBody := testJSONStruct{fieldString, fieldInt}
	reqBodyJSON, _ := json.Marshal(reqBody)

	badlyFormedJSON := strings.Replace(string(reqBodyJSON), "\"value\"", "123", 1)

	wantCode := http.StatusBadRequest
	wantField := "FieldString"
	wantPos := 18

	wantMsg := fmt.Sprintf(
		"Request body contains an invalid value for the %q field (at position %d)",
		wantField,
		wantPos,
	)

	req, _ := http.NewRequest("", "", strings.NewReader(badlyFormedJSON))
	req.Header.Add(consts.ContentType, consts.ApplicationJSON)
	w := httptest.NewRecorder()

	err := validateJSONRequestBody(w, req, &js)

	if err == nil {
		t.Fatal("Error was nil")
	}

	mErr, ok := err.(malformedRequestErr)
	if !ok {
		t.Fatal("Error is not malformedRequest, but it should")
	}

	if mErr.status != wantCode {
		t.Errorf("Wrong http status code, expected %d, but was %d", wantCode, mErr.status)
	}
	if mErr.msg != wantMsg {
		t.Errorf("Expected %s\nbut was %s", wantMsg, mErr.msg)
	}
}

func Test_decodeJSONBodyUnknownField(t *testing.T) {
	var js testJSONStruct
	fieldString, fieldInt := "value", 123
	fieldUnknown := "unknown"

	reqBody := testJSONStruct{fieldString, fieldInt}
	reqBodyJSON, _ := json.Marshal(reqBody)

	withUnknownField := strings.Replace(
		string(reqBodyJSON),
		"}",
		fmt.Sprintf(",\"fieldUnknown\":\"%s\"}", fieldUnknown),
		1,
	)

	wantCode := http.StatusBadRequest
	wantMsg := "Request body contains unknown field \"fieldUnknown\""

	req, _ := http.NewRequest("", "", strings.NewReader(withUnknownField))
	req.Header.Add(consts.ContentType, consts.ApplicationJSON)
	w := httptest.NewRecorder()

	err := validateJSONRequestBody(w, req, &js)

	if err == nil {
		t.Fatal("Error was nil")
	}

	mErr, ok := err.(malformedRequestErr)
	if !ok {
		t.Fatal("Error is not malformedRequest, but it should")
	}

	if mErr.status != wantCode {
		t.Errorf("Wrong http status code, expected %d, but was %d", wantCode, mErr.status)
	}
	if mErr.msg != wantMsg {
		t.Errorf("Expected %s\nbut was %s", wantMsg, mErr.msg)
	}
}

func Test_decodeJSONBodyEmptyBody(t *testing.T) {
	var js testJSONStruct
	req, _ := http.NewRequest("", "", strings.NewReader(""))
	req.Header.Add(consts.ContentType, consts.ApplicationJSON)
	w := httptest.NewRecorder()

	wantCode := http.StatusBadRequest
	wantMsg := "Request body must not be empty"

	err := validateJSONRequestBody(w, req, &js)

	if err == nil {
		t.Fatal("Error was nil")
	}

	mErr, ok := err.(malformedRequestErr)
	if !ok {
		t.Fatal("Error is not malformedRequest, but it should")
	}

	if mErr.status != wantCode {
		t.Errorf("Wrong http status code, expected %d, but was %d", wantCode, mErr.status)
	}
	if mErr.msg != wantMsg {
		t.Errorf("Expected %s\nbut was %s", wantMsg, mErr.msg)
	}
}

func Test_decodeJSONBodyTooLargeBody(t *testing.T) {
	var js testJSONStruct
	fieldString :=
		"toooooooooBiiiiiiiiiiggggggggggggggOooooooooooffffffffAaaaaaaaVaaaaluuuuuuueeeeee"
	fieldInt := 1234567891234567891

	reqBody := testJSONStruct{fieldString, fieldInt}
	reqBodyJSON, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("", "", bytes.NewReader(reqBodyJSON))
	req.Header.Add(consts.ContentType, consts.ApplicationJSON)
	w := httptest.NewRecorder()

	wantCode := http.StatusRequestEntityTooLarge
	wantMsg := fmt.Sprintf("Request body must not be larger than %dB", maxSize)

	err := validateJSONRequestBody(w, req, &js)

	if err == nil {
		t.Fatal("Error was nil")
	}

	mErr, ok := err.(malformedRequestErr)
	if !ok {
		t.Fatal("Error is not malformedRequest, but it should")
	}

	if mErr.status != wantCode {
		t.Errorf("Wrong http status code, expected %d, but was %d", wantCode, mErr.status)
	}
	if mErr.msg != wantMsg {
		t.Errorf("Expected %s\nbut was %s", wantMsg, mErr.msg)
	}
}

func Test_decodeJSONBodyMoreJSONs(t *testing.T) {
	var js testJSONStruct
	fieldString, fieldInt := "value", 123
	reqBody := testJSONStruct{fieldString, fieldInt}
	reqBodyJSON, _ := json.Marshal(reqBody)

	moreJSONs := strings.Repeat(string(reqBodyJSON), 2)

	req, _ := http.NewRequest("", "", strings.NewReader(moreJSONs))
	req.Header.Add(consts.ContentType, consts.ApplicationJSON)
	w := httptest.NewRecorder()

	wantCode := http.StatusBadRequest
	wantMsg := "Request body must only contain a single JSON object"

	err := validateJSONRequestBody(w, req, &js)

	if err == nil {
		t.Fatal("Error was nil")
	}

	mErr, ok := err.(malformedRequestErr)
	if !ok {
		t.Fatal("Error is not malformedRequest, but it should")
	}

	if mErr.status != wantCode {
		t.Errorf("Wrong http status code, expected %d, but was %d", wantCode, mErr.status)
	}
	if mErr.msg != wantMsg {
		t.Errorf("Expected %s\nbut was %s", wantMsg, mErr.msg)
	}
}

func Test_decodeJSONBodyMissingField(t *testing.T) {
	var js testJSONStruct
	fieldString, fieldInt := "value", 123
	reqBody := testJSONStruct{fieldString, fieldInt}
	reqBodyJSON, _ := json.Marshal(reqBody)

	missingOneField := strings.Replace(string(reqBodyJSON), ",\"FieldInt\":123", "", 1)

	req, _ := http.NewRequest("", "", strings.NewReader(missingOneField))
	req.Header.Add(consts.ContentType, consts.ApplicationJSON)
	w := httptest.NewRecorder()

	wantCode := http.StatusBadRequest
	wantMsg := "Request body is not complete"

	err := validateJSONRequestBody(w, req, &js)

	if err == nil {
		t.Fatal("Error was nil")
	}

	mErr, ok := err.(malformedRequestErr)
	if !ok {
		t.Fatal("Error is not malformedRequest, but it should")
	}

	if mErr.status != wantCode {
		t.Errorf("Wrong http status code, expected %d, but was %d", wantCode, mErr.status)
	}
	if mErr.msg != wantMsg {
		t.Errorf("Expected %s\nbut was %s", wantMsg, mErr.msg)
	}
}
