package common

import (
	"encoding/json"
	"io"
	"net/http"
)

type ErrorMessage struct {
	Error   string `json:"error"`
	Status  int    `json:"status"`
	Message string `json:"message"`
}

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

// Writes error to w, sends error status code
func WriteError(w http.ResponseWriter, e error, status int, m ...string) (err error) {
	json := ErrorMessage{
		Error:  e.Error(),
		Status: status,
	}

	if len(m) > 0 {
		json.Message = m[0]
	}

	err = MarshalBody(w, status, json)

	return
}
