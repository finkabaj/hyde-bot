package common

import (
	"errors"
	"net/http"
)

var (
	ErrBadRequest = errors.New("bad request")
	ErrNotFound   = errors.New("not found")
	ErrInternal   = errors.New("internal error")
)

type ErrorResponse struct {
	Error            string            `json:"error,omitempty"`
	Status           int               `json:"status,omitempty"`
	Message          string            `json:"message,omitempty"`
	ValidationErrors map[string]string `json:"validation_errors,omitempty"`
}

type ErrorResponseBuilder interface {
	SetMessage(m string) ErrorResponseBuilder
	SetValidationFields(fields map[string]string) ErrorResponseBuilder
	SetStatus(s int) ErrorResponseBuilder
	Send(w http.ResponseWriter) error
	Get() *ErrorResponse
}

type errorResponseBuilder struct {
	errorResponse *ErrorResponse
}

func NewErrorResponseBuilder(err error) ErrorResponseBuilder {
	return &errorResponseBuilder{
		errorResponse: &ErrorResponse{
			Error: err.Error(),
		},
	}
}

func (e *errorResponseBuilder) SetMessage(m string) ErrorResponseBuilder {
	e.errorResponse.Message = m
	return e
}

func (e *errorResponseBuilder) SetValidationFields(fields map[string]string) ErrorResponseBuilder {
	e.errorResponse.ValidationErrors = fields
	return e
}

func (e *errorResponseBuilder) SetStatus(s int) ErrorResponseBuilder {
	e.errorResponse.Status = s
	return e
}

// Send sends the error response to the client with the status code. If the status code is not set, it will default to 500.
// Returns an error if the response could not be sent.
func (e *errorResponseBuilder) Send(w http.ResponseWriter) (err error) {
	if e.errorResponse.Status == 0 {
		e.errorResponse.Status = http.StatusInternalServerError
	}

	err = MarshalBody(w, e.errorResponse.Status, e.errorResponse)

	return
}

func (e *errorResponseBuilder) Get() *ErrorResponse {
	return e.errorResponse
}
