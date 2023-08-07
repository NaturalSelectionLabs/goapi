package goapi

import (
	"encoding/json"
	"net/http"
	"reflect"
)

func writeResponse(w http.ResponseWriter, data, meta any, err error) {
	if err != nil {
		writeJSON(w, http.StatusBadRequest, Response{
			Error: toError(err),
		})

		return
	}

	writeJSON(w, http.StatusOK, Response{
		Data: data,
		Meta: meta,
	})
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

// Response is a response object that contains the primary data returned by the API.
type Response struct {
	// Data is the primary data returned by the API.
	Data any `json:"data,omitempty"`
	// Meta is a meta object that contains non-standard meta-information,
	// such as pagination information.
	Meta  any    `json:"meta,omitempty"`
	Error *Error `json:"error,omitempty"`
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
