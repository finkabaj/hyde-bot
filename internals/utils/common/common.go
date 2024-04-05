package common

import (
	"encoding/json"
	"io"
	"net/http"
)

func UnmarshalResponse(body io.ReadCloser, v any) (err error) {
	err = json.NewDecoder(body).Decode(v)

	return
}

func MarshalRequest(w http.ResponseWriter, v any) (err error) {
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(v)

	return
}
