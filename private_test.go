package goapi

import (
	"bytes"
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"github.com/ysmood/got"
)

func Test_toValue(t *testing.T) {
	g := got.T(t)

	v, err := toValue(reflect.TypeOf(0), "10")
	g.E(err)

	g.Eq(v.Interface(), 10)
}

func Test_loadURL(t *testing.T) {
	g := got.T(t)

	path, err := newPath("/test/{a}/{d}")
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

	parsed := parseParam(path, reflect.TypeOf(testParams{}))

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

	g.Eq(g.Panic(func() {
		parseParam(path, reflect.TypeOf(struct{}{}))
	}), "expect parameter to be a struct and embedded with "+
		"goapi.InHeader, goapi.InURL, or goapi.InBody, but got: struct {}")
}

func Test_loadURL_err(t *testing.T) {
	g := got.T(t)

	path, err := newPath("/test/{a}")
	g.E(err)

	type testParams struct {
		InURL
	}

	g.Eq(g.Panic(func() {
		parseParam(path, reflect.TypeOf(testParams{}))
	}), "expect to have path parameter for {a} in goapi.testParams")

	type testPath struct {
		InURL
		A int
	}

	parsed := parseParam(path, reflect.TypeOf(testPath{}))

	_, err = parsed.loadURL(url.Values{})
	g.Eq(err.Error(), "missing parameter in request, param: a")

	type testQuery struct {
		InURL
		A int
		B int
	}

	parsed = parseParam(path, reflect.TypeOf(testQuery{}))

	_, err = parsed.loadURL(url.Values{"a": {"1"}})
	g.Eq(err.Error(), "missing parameter in request, param: b")
}

func Test_loadHeader(t *testing.T) {
	g := got.T(t)

	type header struct {
		InHeader
		X_Y int
		Z   string `default:"\"default\""`
	}

	parsed := parseParam(nil, reflect.TypeOf(header{}))

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
		InBody
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	parsed := parseParam(nil, reflect.TypeOf(body{}))

	v, err := parsed.loadBody(bytes.NewBufferString(`{"id": 1, "name": "test"}`))
	g.E(err)

	g.Eq(v.Interface(), body{
		InBody: InBody{},
		ID:     1,
		Name:   "test",
	})
}
