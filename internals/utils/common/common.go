package common

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"reflect"
	"slices"

	"github.com/finkabaj/hyde-bot/internals/utils/rule"
	"github.com/go-playground/validator/v10"
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
	if string(body) == "[]" {
		// If the JSON string is an empty array, set the target to an empty slice
		reflect.ValueOf(v).Elem().Set(reflect.MakeSlice(reflect.TypeOf(v).Elem(), 0, 0))
		return nil
	}

	err = json.Unmarshal(body, v)

	return
}

func ValidateSliceOrStruct(w http.ResponseWriter, validate *validator.Validate, v any) (haveError bool) {
	isSlice := reflect.TypeOf(v).Kind() == reflect.Slice

	var err error

	if isSlice {
		err = validate.Var(v, "required,dive")
	} else {
		err = validate.Struct(v)
	}

	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			// if you see this error that means that it's time to correct validate_json implementation (or you fucked up json)
			SendBadRequestError(w, "Invalid json while validation body")
			return true
		}
		validationErrors := make(map[string]string)
		for _, e := range err.(validator.ValidationErrors) {
			validationErrors[e.Field()] = e.Tag()
		}

		SendValidationError(w, validationErrors)

		return true
	}

	return
}

func EveryFieldValueContains[T any](arr []T, fieldName string, fieldValue interface{}) bool {
	for _, item := range arr {
		v := reflect.ValueOf(item)
		if v.Kind() != reflect.Struct {
			return false
		}

		field := v.FieldByName(fieldName)
		if !field.IsValid() {
			return false
		}

		if field.Interface() != fieldValue {
			return false
		}
	}
	return true
}

func DestructureStructSlice(slice interface{}) [][]any {
	val := reflect.ValueOf(slice)
	if val.Kind() != reflect.Slice {
		panic("Input must be a slice")
	}

	result := make([][]any, val.Len())

	for i := 0; i < val.Len(); i++ {
		structVal := val.Index(i)
		if structVal.Kind() != reflect.Struct {
			panic("Slice elements must be structs")
		}

		fields := make([]any, structVal.NumField())
		for j := 0; j < structVal.NumField(); j++ {
			fields[j] = structVal.Field(j).Interface()
		}

		result[i] = fields
	}

	return result
}

func RemoveDuplicates[T comparable](slice []T) []T {
	seen := make(map[T]bool)
	result := []T{}

	for _, val := range slice {
		if !seen[val] {
			seen[val] = true
			result = append(result, val)
		}
	}

	return result
}

func HaveDuplicatesActions(a [rule.ReactActionCount]rule.ReactAction) bool {
	seen := make(map[rule.ReactAction]bool)

	for _, val := range a {
		if val != 0 && seen[val] {
			return true
		}
		seen[val] = true
	}

	return false
}

func HaveIntersection[T comparable](a, b []T) bool {
	for _, val := range a {
		if slices.Contains(b, val) {
			return true
		}
	}

	return false

}
