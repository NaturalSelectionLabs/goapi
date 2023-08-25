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

type ResponseHeader interface {
	Header() http.Header
}

type ResponseOK struct {
	// Data is the primary data returned by the API.
	Data any `json:"data,omitempty"`
	// Meta is a meta object that contains non-standard meta-information,
	// such as pagination information.
	Meta any `json:"meta,omitempty"`
}

func (r *ResponseOK) StatusCode() int {
	return http.StatusOK
}

type ResponseOKHeader struct {
	ResHeader http.Header `json:"-"`
	Data      any         `json:"data,omitempty"`
	Meta      any         `json:"meta,omitempty"`
}

func (r *ResponseOKHeader) StatusCode() int {
	return http.StatusOK
}

func (r *ResponseOKHeader) Header() http.Header {
	return r.ResHeader
}

type ResponseError struct {
	Error *Error `json:"error,omitempty"`
}

func (r *ResponseError) StatusCode() int {
	return http.StatusBadRequest
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
