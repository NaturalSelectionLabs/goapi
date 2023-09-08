package goapi

import (
	"encoding/json"
	"net/http"
	"reflect"
)

// Response is an interface that represents a response object.
//
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

type parsedRes struct {
	statusCode int
	hasHeader  bool
	hasErr     bool
	hasData    bool
	hasMeta    bool

	typ reflect.Type

	header reflect.Type
	err    reflect.Type
	data   reflect.Type
	meta   reflect.Type
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

func (s *parsedRes) write(path string, w http.ResponseWriter, res reflect.Value) {
	if s.hasHeader {
		h := res.FieldByName("Header")
		for i := 0; i < h.NumField(); i++ {
			f := h.Field(i)
			w.Header().Set(toHeaderName(s.header.Field(i).Name), f.String())
		}
	}

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
	} else if s.hasData {
		data = responseData{
			Data: res.FieldByName("Data").Interface(),
		}
	}

	if s.hasErr || s.hasData {
		b, err := json.Marshal(data)
		if err != nil {
			panic(path + " " + err.Error())
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(s.statusCode)
		_, _ = w.Write(b)
	} else {
		w.WriteHeader(s.statusCode)
	}
}
