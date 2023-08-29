package goapi_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/NaturalSelectionLabs/goapi"
	"github.com/ysmood/got"
)

type res struct {
	goapi.Status200
	Data string
}

type resMeta struct {
	goapi.Status200
	Data string
	Meta string
}

type resErr struct {
	goapi.Status400
	Error goapi.Error
}

type resHeader struct {
	goapi.Status200
	Data   int
	Header struct {
		X_UA string
	}
}

func TestOperation(t *testing.T) {
	g := got.T(t)

	r := goapi.New()

	r.GET("/params", func(params struct {
		goapi.InURL
		A string
	},
	) res {
		return res{Data: params.A}
	})

	r.GET("/meta", func(params struct {
		goapi.InURL
		A string
	},
	) resMeta {
		return resMeta{Data: params.A, Meta: params.A}
	})

	r.GET("/params-time/{t}", func(params struct {
		goapi.InURL
		T time.Time
	},
	) res {
		return res{Data: params.T.String()}
	})

	r.GET("/error-res", func() resErr {
		return resErr{
			Error: goapi.Error{
				Code: "error",
			},
		}
	})

	r.GET("/override-res", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotModified)
	})

	r.GET("/override-header", func() resHeader {
		return resHeader{
			Header: struct{ X_UA string }{X_UA: "test-client"},
		}
	})

	tr := g.Serve()
	tr.Mux.Handle("/", r.Server())

	g.Eq(g.Req("", tr.URL("/params?a=ok")).JSON(), map[string]any{
		"data": "ok",
	})

	g.Eq(g.Req("", tr.URL("/params-time/2023-09-05T14:09:01.123Z")).JSON(), map[string]any{
		"data": "2023-09-05 14:09:01.123 +0000 UTC",
	})

	g.Eq(g.Req("", tr.URL("/error-res")).JSON(), map[string]interface{}{
		"error": map[string]interface{}{
			"code": "error",
		},
	})

	g.Eq(g.Req("", tr.URL("/override-res")).StatusCode, http.StatusNotModified)

	g.Eq(g.Req("", tr.URL("/override-header")).Header.Get("x-ua"), "test-client")
}
