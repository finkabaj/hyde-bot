package middleware

import (
	"fmt"
	"net/http"

	"github.com/finkabaj/hyde-bot/internals/utils/common"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

func ValidateParamId(rule string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := chi.URLParam(r, "id")

			fmt.Println(id)

			if id == "" {
				common.NewErrorResponseBuilder(common.ErrBadRequest).
					SetStatus(http.StatusBadRequest).
					SetMessage("Id is not provided").
					Send(w)
				return
			}

			if err := validate.Var(id, rule); err != nil {
				validateErrors := make(map[string]string)
				vErr := err.(validator.ValidationErrors)[0]
				validateErrors["Id"] = vErr.Error()

				common.NewErrorResponseBuilder(common.ErrValidation).
					SetStatus(http.StatusBadRequest).
					SetValidationFields(validateErrors).
					Send(w)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
