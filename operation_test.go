package goapi_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/NaturalSelectionLabs/goapi"
	"github.com/ysmood/got"
)

type res struct {
	goapi.StatusOK
	Data string
}

type resMeta struct {
	goapi.StatusOK
	Data string
	Meta string
}

type resErr struct {
	goapi.StatusBadRequest
	Error goapi.Error
}

type resHeader struct {
	goapi.StatusOK
	Data   int
	Header struct {
		X_UA string
	}
}

type resEncErr struct {
	goapi.StatusOK
	Data func()
}

type resEmpty struct {
	goapi.StatusOK
}

func TestOperation(t *testing.T) {
	g := got.T(t)
	tr := g.Serve()
	r := goapi.New()

	{ // setup
		r.GET("/query", func(params struct {
			goapi.InURL
			A string
		},
		) res {
			return res{Data: params.A}
		})

		r.GET("/meta", func(params struct {
			goapi.InHeader
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

		r.POST("/req-body", func(params struct {
			goapi.InBody
			A string
		},
		) res {
			return res{Data: params.A}
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

		r.GET("/res-enc-err", func() resEncErr {
			return resEncErr{
				Data: func() {},
			}
		})

		r.GET("/res-empty", func() resEmpty {
			return resEmpty{}
		})

		tr.Mux.Handle("/", r.Server())
	}

	g.Eq(g.Req("", tr.URL("/query?a=ok")).JSON(), map[string]any{
		"data": "ok",
	})

	g.Eq(g.Req("", tr.URL("/query")).StatusCode, http.StatusBadRequest)

	g.Eq(g.Req("", tr.URL("/meta"), http.Header{"a": {"ok"}}).JSON(), map[string]any{
		"data": "ok",
		"meta": "ok",
	})

	g.Eq(g.Req("", tr.URL("/params-time/2023-09-05T14:09:01.123Z")).JSON(), map[string]any{
		"data": "2023-09-05 14:09:01.123 +0000 UTC",
	})

	g.Eq(g.Req(http.MethodPost, tr.URL("/req-body"), `{"a": "ok"}`).JSON(), map[string]any{
		"data": "ok",
	})

	g.Eq(g.Req("", tr.URL("/error-res")).JSON(), map[string]interface{}{
		"error": map[string]interface{}{
			"code": "error",
		},
	})

	g.Eq(g.Req("", tr.URL("/override-res")).StatusCode, http.StatusNotModified)

	g.Eq(g.Req("", tr.URL("/override-header")).Header.Get("x-ua"), "test-client")

	g.Eq(g.Panic(func() {
		r.GET("/[", func() res { return res{} })
	}).(error).Error(), "error parsing regexp: missing closing ]: `[$`")

	g.Eq(g.Panic(func() {
		r.GET("/", 10)
	}), "handler must be a function")

	g.Eq(g.Panic(func() {
		r.GET("/", func() {})
	}), "handler must return a single value")

	g.Eq(g.Req("", tr.URL("/res-enc-err")).StatusCode, http.StatusInternalServerError)

	g.Eq(g.Req("", tr.URL("/res-empty")).String(), "")
}
