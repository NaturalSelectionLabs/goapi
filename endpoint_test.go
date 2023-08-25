package goapi_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/NaturalSelectionLabs/goapi"
	"github.com/ysmood/got"
)

func TestEndpoint(t *testing.T) {
	g := got.T(t)

	r := goapi.New()

	r.GET("/no-content", func() {})

	r.GET("/data", func() string {
		return "ok"
	})

	r.GET("/data-meta", func() (string, string) {
		return "hello", "world"
	})

	r.GET("/internal-err", func() error {
		return fmt.Errorf("error")
	})

	r.GET("/bad-request", func() error {
		return &goapi.Error{
			Code:    "bad-request",
			Message: "bad request",
		}
	})

	r.GET("/error-res", func() goapi.Response {
		return &goapi.ResponseBadRequest{
			Error: &goapi.Error{
				Code: "error",
			},
		}
	})

	r.GET("/override-res", func(w http.ResponseWriter) {
		w.WriteHeader(http.StatusNotModified)
	})

	r.GET("/override-header", func() goapi.Response {
		return &goapi.ResponseOKHeader{
			ResHeader: http.Header{
				"x-test": []string{"ok"},
			},
		}
	})

	g.Eq(g.Panic(func() {
		r.GET("/bad-params/{p}", func(p string) {})
	}), "expect handler arguments must be http.ResponseWriter, *http.Request, or pointer to a struct, but got: string")

	tr := g.Serve()
	tr.Mux.Handle("/", r)

	g.Eq(g.Req("", tr.URL("/no-content")).StatusCode, http.StatusNoContent)

	g.Eq(g.Req("", tr.URL("/data")).JSON(), map[string]interface{}{
		"data": "ok",
	})

	g.Eq(g.Req("", tr.URL("/data-meta")).JSON(), map[string]interface{}{
		"data": "hello",
		"meta": "world",
	})

	g.Eq(g.Req("", tr.URL("/internal-err")).StatusCode, http.StatusInternalServerError)

	g.Eq(g.Req("", tr.URL("/bad-request")).StatusCode, http.StatusBadRequest)

	g.Eq(g.Req("", tr.URL("/error-res")).StatusCode, http.StatusBadRequest)
	g.Eq(g.Req("", tr.URL("/error-res")).JSON(), map[string]interface{}{
		"error": map[string]interface{}{
			"code": "error",
		},
	})

	g.Eq(g.Req("", tr.URL("/override-res")).StatusCode, http.StatusNotModified)

	g.Eq(g.Req("", tr.URL("/override-header")).Header.Get("x-test"), "ok")
}
