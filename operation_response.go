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

type FormatResponse func(ResponseFormat) any

type ResponseFormat interface {
	format()
}

type ResponseFormatErr struct {
	Error any `json:"error"`
}

func (ResponseFormatErr) format() {}

type ResponseFormatMeta struct {
	Data any `json:"data"`
	Meta any `json:"meta"`
}

func (ResponseFormatMeta) format() {}

type ResponseFormatData struct {
	Data any `json:"data"`
}

func (ResponseFormatData) format() {}

type parsedRes struct {
	operation  *Operation
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

func (op *Operation) parseResponse(t reflect.Type) *parsedRes {
	if !t.Implements(tResponse) {
		panic("handler must return a goapi.Response")
	}

	res := &parsedRes{operation: op, typ: t}

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

func (s *parsedRes) write(w http.ResponseWriter, res reflect.Value) {
	if s.hasHeader {
		h := res.FieldByName("Header")
		for i := 0; i < h.NumField(); i++ {
			f := h.Field(i)
			w.Header().Set(toHeaderName(s.header.Field(i).Name), f.String())
		}
	}

	var format ResponseFormat

	if s.hasErr { //nolint: gocritic
		format = ResponseFormatErr{
			Error: res.FieldByName("Error").Interface(),
		}
	} else if s.hasMeta {
		format = ResponseFormatMeta{
			Data: res.FieldByName("Data").Interface(),
			Meta: res.FieldByName("Meta").Interface(),
		}
	} else if s.hasData {
		format = ResponseFormatData{
			Data: res.FieldByName("Data").Interface(),
		}
	}

	if s.hasErr || s.hasData {
		b, err := json.Marshal(s.operation.group.router.FormatResponse(format))
		if err != nil {
			panic(s.operation.path.path + " " + err.Error())
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(s.statusCode)
		_, _ = w.Write(b)
	} else {
		w.WriteHeader(s.statusCode)
	}
}
