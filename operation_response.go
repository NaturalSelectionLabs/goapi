package goapi

import (
	"encoding/json"
	"io"
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

// DataStream is a flag for binary response body.
// When Data field in the response struct is of this type,
// the response body will be written directly to the [http.ResponseWriter].
type DataStream io.Reader

var tDataStream = reflect.TypeOf(new(DataStream)).Elem()

type parsedRes struct {
	operation  *Operation
	statusCode int
	hasHeader  bool
	hasErr     bool
	hasData    bool
	hasMeta    bool
	isDirect   bool

	isStream    bool
	contentType string

	typ reflect.Type

	header reflect.Type
	err    reflect.Type
	data   reflect.Type
	meta   reflect.Type
}

const (
	// TagResponse is one tag name for Data field.
	TagResponse = "response"
	// TagResponseDirect is one value for [TagResponse].
	TagResponseDirect = "direct"
)

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

	res.contentType = getContentType(t, "")

	if err, has := t.FieldByName("Error"); has {
		res.hasErr = true
		res.err = err.Type
	}

	if f, has := t.FieldByName("Data"); has {
		if res.hasErr {
			panic("response Data field should not exist when Error field exists")
		}

		if f.Type == tDataStream {
			res.isStream = true
		}

		if f.Tag.Get(TagResponse) == TagResponseDirect {
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

		if res.isStream {
			panic("response Meta field cannot exist when Data field is goapi.DataStream")
		}

		res.hasMeta = true
		res.meta = f.Type
	}

	return res
}

func (s *parsedRes) write(w http.ResponseWriter, res reflect.Value) {
	if s.contentType != "" {
		w.Header().Set("Content-Type", s.contentType)
	}

	if s.hasHeader {
		h := res.FieldByName("Header")
		for i := 0; i < h.NumField(); i++ {
			f := h.Field(i)
			w.Header().Set(toHeaderName(s.header.Field(i).Name), f.String())
		}
	}

	if s.isStream {
		if w.Header().Get("Content-Type") == "" {
			w.Header().Set("Content-Type", "application/octet-stream")
		}

		data := res.FieldByName("Data").Interface()

		w.WriteHeader(s.statusCode)
		_, _ = io.Copy(w, data.(DataStream))

		if closer, ok := data.(io.Closer); ok {
			_ = closer.Close()
		}

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
