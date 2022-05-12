package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	CONTENT_TYPE     = "Content-Type"
	APPLICATION_JSON = "application/json"

	MB = 1048576
)

type malformedRequest struct {
	status int
	msg    string
}

func (mr *malformedRequest) Error() string {
	return mr.msg
}

func decodeJSONBody[T any](w http.ResponseWriter, r *http.Request, dest T) error {

	if ct := r.Header.Get(CONTENT_TYPE); ct != APPLICATION_JSON {
		responseMsg := "Content-Type header is not application/json"
		return &malformedRequest{status: http.StatusUnsupportedMediaType, msg: responseMsg}
	}

	r.Body = http.MaxBytesReader(w, r.Body, MB)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&dest)

	if err != nil {

		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			responseMsg := fmt.Sprintf(
				"Request body contains badly-formed JSON (at position %d)",
				syntaxError.Offset,
			)
			return &malformedRequest{status: http.StatusBadRequest, msg: responseMsg}

		case errors.Is(err, io.ErrUnexpectedEOF):
			responseMsg := fmt.Sprintf("Request body contains badly-formed JSON")
			return &malformedRequest{status: http.StatusBadRequest, msg: responseMsg}

		case errors.As(err, &unmarshalTypeError):
			responseMsg := fmt.Sprintf(
				"Request body contains an invalid value for the %q field (at position %d)",
				unmarshalTypeError.Field,
				unmarshalTypeError.Offset,
			)
			return &malformedRequest{status: http.StatusBadRequest, msg: responseMsg}

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			responseMsg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			return &malformedRequest{status: http.StatusBadRequest, msg: responseMsg}

		case errors.Is(err, io.EOF):
			responseMsg := "Request body must not be empty"
			return &malformedRequest{status: http.StatusBadRequest, msg: responseMsg}

		case err.Error() == "http: request body too large":
			responseMsg := "Request body must not be larger than 1MB"
			return &malformedRequest{status: http.StatusRequestEntityTooLarge, msg: responseMsg}

		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		responseMsg := "Request body must only contain a single JSON object"
		return &malformedRequest{status: http.StatusBadRequest, msg: responseMsg}
	}

	return nil
}
