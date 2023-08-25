package goapi_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/NaturalSelectionLabs/goapi"
	"github.com/ysmood/got"
)

func TestEndpoint(t *testing.T) {
	g := got.T(t)

	r := goapi.New()

	r.GET("/query-int", func(params *struct{ ID int }) int { return params.ID })

	r.GET("/query-float", func(params *struct{ ID float64 }) float64 { return params.ID })

	r.GET("/query-ptr", func(params *struct{ ID *int }) int { return *params.ID })

	r.GET("/query-arr", func(params *struct{ ID []int }) int { return params.ID[1] })

	r.GET("/query-decoder", func(params *struct{ T Time }) int64 { return params.T.t.Unix() })

	r.GET("/query-decoder/{t}", func(params *struct{ T *Time }) int64 { return params.T.t.Unix() })

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

	r.GET("/override-header", func(r *http.Request) goapi.Response {
		return &goapi.ResponseOKHeader{
			ResHeader: http.Header{
				"x-ua": []string{r.UserAgent()},
			},
		}
	})

	g.Eq(g.Panic(func() {
		r.GET("/bad-params/{p}", func(p string) {})
	}), "expect handler arguments must be http.ResponseWriter, *http.Request, or pointer to a struct, but got: string")

	g.Eq(g.Panic(func() {
		r.GET("/bad-params/[a/", func() {})
	}), "expect path matches the openapi path format, but got: /bad-params/[a/")

	tr := g.Serve()
	tr.Mux.Handle("/", r)

	g.Eq(g.Req("", tr.URL("/query-int?id=1")).JSON(), map[string]interface{}{
		"data": 1.0,
	})

	g.Eq(g.Req("", tr.URL("/query-float?id=1.2")).JSON(), map[string]interface{}{
		"data": 1.2,
	})

	g.Eq(g.Req("", tr.URL("/query-ptr?id=1")).JSON(), map[string]interface{}{
		"data": 1.0,
	})

	g.Eq(g.Req("", tr.URL("/query-arr?id=1&id=2&id=3")).JSON(), map[string]interface{}{
		"data": 2.0,
	})

	g.Eq(g.Req("", tr.URL("/query-decoder?t=2023-01-02")).JSON(), map[string]interface{}{
		"data": 1672617600.0,
	})

	g.Eq(g.Req("", tr.URL("/query-decoder/2023-01-02")).JSON(), map[string]interface{}{
		"data": 1672617600.0,
	})

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

	g.Has(g.Req("", tr.URL("/override-header")).Header.Get("x-ua"), "Go-http-client")
}

type Time struct {
	t time.Time
}

func (t *Time) DecodeParam(v []string) {
	tt, err := time.Parse(time.DateOnly, v[0])
	if err != nil {
		panic(err)
	}

	t.t = tt
}

func TestSetTags(t *testing.T) {
	g := got.T(t)

	r := goapi.New()
	e := r.GET("/", func() {})
	e.SetTags(goapi.NewTag("tag"))

	g.Eq(e.Tags()[0].String(), "tag")

	g.Eq(1, 1)
}
