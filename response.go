package goapi

import (
	"encoding/json"
	"net/http"
	"reflect"
)

// Response is a response object that contains the primary data returned by the API.
type Response interface {
	// StatusCode is the HTTP status code of the response.
	// It should always return the same value.
	StatusCode() int
}

// ResponseHeader is a [Response] object that implements extra Header method.
type ResponseHeader interface {
	Header() http.Header
}

// ErrorCode is an interface to tell if a error is a managed error.
type ErrorCode interface {
	ErrorCode() string
}

type ResponseOK struct {
	// Data is the primary data returned by the API.
	Data any `json:"data"`
	// Meta is a meta object that contains non-standard meta-information,
	// such as pagination information.
	Meta any `json:"meta,omitempty"`
}

func (r *ResponseOK) StatusCode() int {
	return http.StatusOK
}

// ResponseOKHeader is similar with [ResponseOK] but with an extra ResHeader field.
type ResponseOKHeader struct {
	ResHeader http.Header `json:"-"`
	Data      any         `json:"data"`
	Meta      any         `json:"meta,omitempty"`
}

func (r *ResponseOKHeader) StatusCode() int {
	return http.StatusOK
}

func (r *ResponseOKHeader) Header() http.Header {
	return r.ResHeader
}

type ResponseBadRequest struct {
	Error *Error `json:"error"`
}

func (r *ResponseBadRequest) StatusCode() int {
	return http.StatusBadRequest
}

type ResponseNotFound struct {
	Error *Error `json:"error,omitempty"`
}

func (r *ResponseNotFound) StatusCode() int {
	return http.StatusNotFound
}

type ResponseInternalServerError struct {
	Error *Error `json:"error,omitempty"`
}

func (r *ResponseInternalServerError) StatusCode() int {
	return http.StatusInternalServerError
}

// Error is an error object that contains information about a failed request.
// Reference: https://github.com/microsoft/api-guidelines/blob/vNext/Guidelines.md#error--object
type Error struct {
	// Code is a machine-readable error code.
	Code string `json:"code,omitempty"`
	// Message is a human-readable error message.
	Message string `json:"message,omitempty"`
	// Target is a human-readable description of the target of the error.
	Target string `json:"target,omitempty"`
	// Details is an array of structured error details objects.
	Details []Error `json:"details,omitempty"`
	// InnerError is a generic error object that is used by the service developer for debugging.
	InnerError any `json:"innererror,omitempty"`
}

func (e *Error) ErrorCode() string {
	return e.Code
}

func (e *Error) Error() string {
	return e.Message
}

func toError(err error) *Error {
	if err == nil {
		return nil
	}

	if e, ok := err.(*Error); ok { //nolint: errorlint
		return e
	}

	return &Error{
		Code:    reflect.TypeOf(err).String(),
		Message: err.Error(),
	}
}

func writeResponse(w http.ResponseWriter, res Response) {
	w.Header().Set("Content-Type", "application/json")

	if resHeader, ok := res.(ResponseHeader); ok {
		for k, v := range resHeader.Header() {
			w.Header()[k] = v
		}
	}

	w.WriteHeader(res.StatusCode())
	_ = json.NewEncoder(w).Encode(res)
}
