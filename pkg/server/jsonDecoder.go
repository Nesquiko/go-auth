package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

const (
	// maxSize is a maximal size, in Bytes, of a JSON reques body.
	maxSize = 128
)

// malformedRequestErr represents a error caused by a malformed JSON request.
// Status represents what http status code should be returned to user, and
// message is a description of what was wrong with the JSON request.
type malformedRequestErr struct {
	status int
	msg    string
}

// Error returns a string representation of malformedRequest error.
func (mr malformedRequestErr) Error() string {
	return mr.msg
}

// validateJSONRequest validates if the JSON request is correct form. If something
// is not valid, a malformedResultErr error is returned with details about the
// error that occured.
func validateJSONRequest[T any](w http.ResponseWriter, r *http.Request, dest T) error {

	if ct := r.Header.Get(contentType); ct != applicationJSON {
		responseMsg := "Content-Type header is not application/json"
		return malformedRequestErr{status: http.StatusUnsupportedMediaType, msg: responseMsg}
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxSize)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&dest)

	if err != nil {
		return analyzeError(err)
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		responseMsg := "Request body must only contain a single JSON object"
		return malformedRequestErr{status: http.StatusBadRequest, msg: responseMsg}
	}

	if err = validateDest(dest); err != nil {
		return err
	}

	return nil
}

// validateDest validates if all required fields were populated from JSON
// request. If not, appropriate malformedRequestErr is returned.
func validateDest[T any](dest T) error {
	v := validator.New()
	valErr := v.Struct(dest)
	if valErr != nil {
		responseMsg := "Request body is not complete"
		return malformedRequestErr{status: http.StatusBadRequest, msg: responseMsg}
	}

	return nil
}

// analyzeError tries to specify what the param err is about and then returns
// appropriate malformedRequestErr error.
func analyzeError(err error) malformedRequestErr {

	var syntaxError *json.SyntaxError
	var unmarshalTypeError *json.UnmarshalTypeError

	switch {
	case errors.As(err, &syntaxError):
		responseMsg := fmt.Sprintf(
			"Request body contains badly-formed JSON (at position %d)",
			syntaxError.Offset,
		)
		return malformedRequestErr{status: http.StatusBadRequest, msg: responseMsg}

	case errors.Is(err, io.ErrUnexpectedEOF):
		responseMsg := fmt.Sprintf("Request body contains badly-formed JSON")
		return malformedRequestErr{status: http.StatusBadRequest, msg: responseMsg}

	case errors.As(err, &unmarshalTypeError):
		responseMsg := fmt.Sprintf(
			"Request body contains an invalid value for the %q field (at position %d)",
			unmarshalTypeError.Field,
			unmarshalTypeError.Offset,
		)
		return malformedRequestErr{status: http.StatusBadRequest, msg: responseMsg}

	case strings.HasPrefix(err.Error(), "json: unknown field "):
		fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
		responseMsg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
		return malformedRequestErr{status: http.StatusBadRequest, msg: responseMsg}

	case errors.Is(err, io.EOF):
		responseMsg := "Request body must not be empty"
		return malformedRequestErr{status: http.StatusBadRequest, msg: responseMsg}

	case err.Error() == "http: request body too large":
		responseMsg := fmt.Sprintf("Request body must not be larger than %dB", maxSize)
		return malformedRequestErr{status: http.StatusRequestEntityTooLarge, msg: responseMsg}

	default:
		return malformedRequestErr{status: http.StatusBadRequest, msg: "unknown error"}
	}
}
