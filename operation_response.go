package goapi

import (
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/NaturalSelectionLabs/goapi/lib/openapi"
)

// Response is an interface that represents a response object.
//
//go:generate go run ./lib/gen-status-code
type Response interface {
	statusCode() int
}

var tResponse = reflect.TypeOf((*Response)(nil)).Elem()

// DataBinary is a flag for binary response body.
// When Data field in the response struct is of this type,
// the response body will be written directly to the [http.ResponseWriter].
type DataBinary []byte

var tDataBinary = reflect.TypeOf(DataBinary{})

type parsedRes struct {
	operation  *Operation
	statusCode int
	hasHeader  bool
	hasErr     bool
	hasData    bool
	hasMeta    bool
	isDirect   bool
	isBinary   bool

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

		if f.Type == tDataBinary {
			res.isBinary = true
		}

		if f.Tag.Get("response") == "direct" {
			res.isDirect = true
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

		if res.isBinary {
			panic("response Meta field cannot exist when Data field is goapi.DataBinary")
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

	if s.isBinary {
		w.Header().Set("Content-Type", "application/octet-stream")
		w.WriteHeader(s.statusCode)
		_, _ = w.Write(res.FieldByName("Data").Interface().(DataBinary))

		return
	}

	var data any

	if s.isDirect {
		data = res.FieldByName("Data").Interface()
	} else {
		var format openapi.ResponseFormat

		if s.hasErr { //nolint: gocritic
			format = openapi.ResponseFormatErr{
				Error: res.FieldByName("Error").Interface(),
			}
		} else if s.hasMeta {
			format = openapi.ResponseFormatMeta{
				Data: res.FieldByName("Data").Interface(),
				Meta: res.FieldByName("Meta").Interface(),
			}
		} else if s.hasData {
			format = openapi.ResponseFormatData{
				Data: res.FieldByName("Data").Interface(),
			}
		}

		data = format
	}

	if s.hasErr || s.hasData {
		b, err := json.Marshal(data)
		if err != nil {
			panic(s.operation.path.path + " " + err.Error())
		}

		setJSONHeader(w)
		w.WriteHeader(s.statusCode)
		_, _ = w.Write(b)
	} else {
		w.WriteHeader(s.statusCode)
	}
}

func setJSONHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
}
