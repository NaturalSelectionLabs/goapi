package goapi

import (
	"encoding/json"
	"net/http"
	"reflect"
)

//go:generate go run ./lib/gen-status-code
type Response interface {
	statusCode() int
}

var tResponse = reflect.TypeOf((*Response)(nil)).Elem()

type responseErr struct {
	Error any `json:"error"`
}

type responseMeta struct {
	Data any `json:"data"`
	Meta any `json:"meta"`
}

type responseData struct {
	Data any `json:"data"`
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

func writeResErr(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)

	err := json.NewEncoder(w).Encode(Error{Message: msg})
	if err != nil {
		panic(err)
	}
}

type parsedRes struct {
	statusCode int
	hasHeader  bool
	hasErr     bool
	hasData    bool
	hasMeta    bool

	typ reflect.Type

	header reflect.Type
	data   reflect.Type
	meta   reflect.Type
	err    reflect.Type
}

func parseResponse(t reflect.Type) *parsedRes {
	if !t.Implements(tResponse) {
		panic("handler must return a goapi.Response")
	}

	res := &parsedRes{typ: t}

	res.statusCode = reflect.New(t).Elem().Interface().(Response).statusCode()

	if header, has := t.FieldByName("Header"); has {
		res.hasHeader = true
		res.header = header.Type
	}

	if err, has := t.FieldByName("Error"); has {
		res.hasErr = true
		res.err = err.Type
	}

	if f, has := t.FieldByName("Data"); has {
		if res.hasErr {
			panic("response Data field should not exist when Error field exists")
		}

		res.hasData = true
		res.data = f.Type
	}

	if !res.hasData && !res.hasErr {
		panic("response must have either Data or Error field")
	}

	if f, has := t.FieldByName("Meta"); has {
		if res.hasErr {
			panic("response Meta field should not exist when Error field exists")
		}

		if !res.hasData {
			panic("response Meta field requires Data field")
		}

		res.hasMeta = true
		res.meta = f.Type
	}

	return res
}

func (s *parsedRes) write(w http.ResponseWriter, res reflect.Value) {
	if s.hasErr || s.hasData {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
	}

	if s.hasHeader {
		h := res.FieldByName("Header")
		for i := 0; i < h.NumField(); i++ {
			f := h.Field(i)
			w.Header().Set(toHeaderName(s.header.Field(i).Name), f.String())
		}
	}

	w.WriteHeader(s.statusCode)

	var data any

	if s.hasErr { //nolint: gocritic
		data = responseErr{
			Error: res.FieldByName("Error").Interface(),
		}
	} else if s.hasMeta {
		data = responseMeta{
			Data: res.FieldByName("Data").Interface(),
			Meta: res.FieldByName("Meta").Interface(),
		}
	} else {
		data = responseData{
			Data: res.FieldByName("Data").Interface(),
		}
	}

	if s.hasErr || s.hasData {
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			panic(err)
		}
	}
}
