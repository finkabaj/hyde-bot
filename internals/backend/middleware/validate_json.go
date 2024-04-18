package middleware

import (
	"context"
	"io"
	"net/http"
	"reflect"
	"strings"

	"github.com/finkabaj/hyde-bot/internals/utils/common"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate
var ValidateJsonCtxKey = "json"

func init() {
	validate = validator.New()
	validate.RegisterTagNameFunc(func(f reflect.StructField) string {
		name := strings.SplitN(f.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})
}

// ValidateJson is a middleware that validates the json body of a request and respond
// with an error if json body is empty or if the json body is invalid.
// To use it, you must pass a struct to generic type T that has json and validate tags and call it as a function.
func ValidateJson[T any]() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body := r.Body
			defer body.Close()

			var v T

			if err := common.UnmarshalBody(body, &v); err == io.EOF {
				common.NewErrorResponseBuilder(common.ErrEmptyBody).
					SetStatus(http.StatusBadRequest).
					SetMessage("empty request body").
					Send(w)
				return
			}

			// XXX  maybe need to update in future.
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
					common.NewErrorResponseBuilder(common.ErrInternal).
						SetStatus(http.StatusInternalServerError).
						SetMessage("incorrect json type for validation").
						Send(w)
					return
				}
				validationErrors := make(map[string]string)
				for _, e := range err.(validator.ValidationErrors) {
					validationErrors[e.Field()] = e.Tag()
				}

				common.NewErrorResponseBuilder(common.ErrValidation).
					SetStatus(http.StatusBadRequest).
					SetValidationFields(validationErrors).
					Send(w)

				return
			}

			ctx := context.WithValue(r.Context(), ValidateJsonCtxKey, v)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
