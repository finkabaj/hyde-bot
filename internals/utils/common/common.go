package common

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

var (
	ErrEmptyBody  = errors.New("empty request body")
	ErrValidation = errors.New("validation error")
)

// Reads json body to v
func UnmarshalBody(body io.ReadCloser, v any) (err error) {
	err = json.NewDecoder(body).Decode(v)

	return
}

// Writes json body to w, sends status code
func MarshalBody(w http.ResponseWriter, status int, v any) (err error) {
	w.Header().Set("Content-Type", "application/json; charset=utf8")
	w.WriteHeader(status)
	err = json.NewEncoder(w).Encode(v)

	return
}
