package goapi

import (
	"bytes"
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"github.com/NaturalSelectionLabs/jschema"
	"github.com/ysmood/got"
)

func Test_paramsInGuard(t *testing.T) {
	g := got.T(t)

	_ = InHeader{}.inHeader()

	_ = InURL{}.inURL()

	g.Eq(1, 1)
}

func Test_toValue(t *testing.T) {
	g := got.T(t)

	v, err := toValue(reflect.TypeOf(0), "10")
	g.E(err)

	g.Eq(v.Interface(), 10)
}

func Test_loadURL(t *testing.T) {
	g := got.T(t)

	path, err := newPath("/test/{a}/{d}", false)
	g.E(err)

	type testParams struct {
		InURL
		A int
		B float64
		C *string
		D bool
		E []int
		F []int
	}

	s := jschema.New("")

	parsed := parseParam(s, path, reflect.TypeOf(testParams{}))

	v, err := parsed.loadURL(url.Values{
		"a": []string{"10"},
		"b": []string{"10.1"},
		"c": []string{"ok"},
		"d": []string{"true"},
		"e": []string{"1", "2"},
	})
	g.E(err)

	g.Eq(v.Interface(), testParams{
		InURL: InURL{},
		A:     10,
		B:     10.1,
		C:     strPtr("ok"),
		D:     true,
		E: []int /* len=2 cap=2 */ {
			1,
			2,
		},
		F: nil,
	})
}

func Test_loadURL_any(t *testing.T) {
	g := got.T(t)

	type params struct {
		InURL

		Path string `path:"*"`
	}

	path, err := newPath("/test/*", false)
	g.E(err)

	s := jschema.New("")

	parsed := parseParam(s, path, reflect.TypeOf(params{}))

	v, err := parsed.loadURL(url.Values{
		"*": []string{"a/b/c"},
	})
	g.E(err)

	g.Eq(v.Interface(), params{
		InURL: InURL{},
		Path:  "a/b/c",
	})
}

func Test_loadURL_nil(t *testing.T) {
	g := got.T(t)

	path, err := newPath("/test", false)
	g.E(err)

	type testParams struct {
		InURL
		A *string
	}

	s := jschema.New("")

	parsed := parseParam(s, path, reflect.TypeOf(testParams{}))

	v, err := parsed.loadURL(url.Values{})
	g.E(err)

	g.Eq(v.Interface(), testParams{
		A: nil,
	})
}

func Test_loadURL_err(t *testing.T) {
	g := got.T(t)

	path, err := newPath("/test/{a}", false)
	g.E(err)

	type testParams struct {
		InURL
	}

	s := jschema.New("")

	g.Eq(g.Panic(func() {
		parseParam(s, path, reflect.TypeOf(testParams{}))
	}), "expect to have path parameter for {a} in goapi.testParams")

	type testPath struct {
		InURL
		A int
	}

	parsed := parseParam(s, path, reflect.TypeOf(testPath{}))

	_, err = parsed.loadURL(url.Values{})
	g.Eq(err.Error(), "missing url path param `a`")

	type testQuery struct {
		InURL
		A int
		B int
		C []int
	}

	parsed = parseParam(s, path, reflect.TypeOf(testQuery{}))

	_, err = parsed.loadURL(url.Values{"a": {"1"}})
	g.Eq(err.Error(), "missing url query param `b`")

	_, err = parsed.loadURL(url.Values{"a": {"true"}})
	g.Eq(err.Error(), "failed to parse url path param `a`: can't parse `true` to expected value, "+
		"json: cannot unmarshal bool into Go value of type int")

	_, err = parsed.loadURL(url.Values{"a": {"1"}, "b": {"2"}, "c": {"true"}})
	g.Eq(err.Error(), "failed to parse url param `c`: can't parse `true` to expected value, "+
		"json: cannot unmarshal bool into Go value of type int")

	g.Eq(g.Panic(func() {
		parseParam(s, path, reflect.TypeOf(struct {
			InURL
			A []int
		}{}))
	}), "path parameter cannot be an slice, param: A")

	g.Eq(g.Panic(func() {
		parseParam(s, path, reflect.TypeOf(struct {
			InURL
			A *int
		}{}))
	}), "path parameter cannot be optional, param: A")

	g.Eq(g.Panic(func() {
		parseParam(s, path, reflect.TypeOf(struct {
			InURL
			A int `default:"1"`
		}{}))
	}), "path parameter cannot have tag `default`, param: A")
}

func Test_loadHeader(t *testing.T) {
	g := got.T(t)

	type header struct {
		InHeader
		X_Y int
		Z   string `default:"default"`
	}

	s := jschema.New("")

	parsed := parseParam(s, nil, reflect.TypeOf(header{}))

	v, err := parsed.loadHeader(http.Header{
		"X-Y": []string{"10"},
	})
	g.E(err)

	g.Eq(v.Interface(), header{
		InHeader: InHeader{},
		X_Y:      10,
		Z:        "default",
	})
}

func strPtr(s string) *string {
	return &s
}

func Test_loadBody(t *testing.T) {
	g := got.T(t)

	type body struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	s := jschema.New("")

	parsed := parseParam(s, nil, reflect.TypeOf(body{}))

	v, err := parsed.loadBody(bytes.NewBufferString(`{"id": 1, "name": "test"}`))
	g.E(err)

	g.Eq(v.Interface(), body{
		ID:   1,
		Name: "test",
	})

	_, err = parsed.loadBody(bytes.NewBufferString(`{`))
	g.Eq(err.Error(), "failed to parse json body: unexpected EOF")
}

func Test_parseResponse_err(t *testing.T) {
	g := got.T(t)
	op := &Operation{}

	g.Eq(g.Panic(func() {
		op.parseResponse(reflect.TypeOf(struct{}{}))
	}), "handler must return a goapi.Response")

	g.Eq(g.Panic(func() {
		op.parseResponse(reflect.TypeOf(struct {
			StatusOK
			Data  int
			Error int
		}{}))
	}), "response Data field should not exist when Error field exists")

	g.Eq(g.Panic(func() {
		op.parseResponse(reflect.TypeOf(struct {
			StatusOK
			Meta  int
			Error int
		}{}))
	}), "response Meta field should not exist when Error field exists")

	g.Eq(g.Panic(func() {
		op.parseResponse(reflect.TypeOf(struct {
			StatusOK
			Meta int
		}{}))
	}), "response Meta field requires Data field")

	g.Eq(g.Panic(func() {
		op.parseResponse(reflect.TypeOf(struct {
			StatusOK
			Data DataStream
			Meta int
		}{}))
	}), "response Meta field cannot exist when Data field is goapi.DataStream")
}

func Test_default_arr(t *testing.T) {
	g := got.T(t)

	type params struct {
		InURL
		IDS []int `default:"[1, 2]"`
	}

	path, err := newPath("/test", false)
	g.E(err)

	s := jschema.New("")

	parsed := parseParam(s, path, reflect.TypeOf(params{}))

	v, err := parsed.loadURL(url.Values{})
	g.E(err)

	g.Eq(v.Interface(), params{
		IDS: []int{1, 2},
	})
}

type myID struct {
}

func (myID) IsFormat(input interface{}) bool {
	return input == "ok"
}

func Test_custom_checker(t *testing.T) {
	g := got.T(t)

	r := New()
	r.Router().AddFormatChecker("my-id", myID{})

	type params struct {
		InURL
		ID string `format:"my-id"`
	}

	path, err := newPath("/test", false)
	g.E(err)

	s := jschema.New("")

	parsed := parseParam(s, path, reflect.TypeOf(params{}))

	_, err = parsed.loadURL(url.Values{"id": {"ok"}})
	g.Nil(err)

	_, err = parsed.loadURL(url.Values{"id": {"no"}})
	g.Eq(err.Error(), "param `id` is invalid: [(root): Does not match format 'my-id']")
}

func Test_validation(t *testing.T) {
	g := got.T(t)

	type A struct {
		ID string `minLen:"5"`
	}

	path, err := newPath("/test", false)
	g.E(err)

	s := jschema.New("")

	{
		parsed := parseParam(s, path, reflect.TypeOf(A{}))

		_, err := parsed.loadBody(bytes.NewBufferString(`{"id": "ok"}`))
		g.Eq(err.Error(), "request body is invalid: [ID: String length must be greater than or equal to 5]")
	}

	{
		type B struct {
			InURL
			ID int `min:"1" max:"10"`
		}

		parsed := parseParam(s, path, reflect.TypeOf(B{}))
		_, err := parsed.loadURL(url.Values{"id": {"0"}})
		g.Eq(err.Error(), "param `id` is invalid: [(root): Must be greater than or equal to 1]")

		_, err = parsed.loadURL(url.Values{"id": {"20"}})
		g.Eq(err.Error(), "param `id` is invalid: [(root): Must be less than or equal to 10]")
	}
}
