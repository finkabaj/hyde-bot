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

type OkResponse struct {
	Message string `json:"message"`
}

func GetApiUrl(host, port, path string) string {
	return "http://" + host + ":" + port + path
}

// Reads json body to v. Body is ReadCloser
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

// Use this function if you have UnmarshalJSON method in your struct
func UnmarshalBodyBytes(body []byte, v any) (err error) {
	err = json.Unmarshal(body, v)

	return
}
