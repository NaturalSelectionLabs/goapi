package goapi_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/NaturalSelectionLabs/goapi"
	"github.com/NaturalSelectionLabs/goapi/lib/middlewares"
	"github.com/ysmood/got"
)

func TestMiddleware(t *testing.T) {
	g := got.T(t)

	r := goapi.NewRouter()

	r.Use(middlewares.Func(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, rq *http.Request) {
			rq = rq.WithContext(context.WithValue(rq.Context(), "middleware01", "ok")) //nolint: staticcheck
			h.ServeHTTP(w, rq)
		})
	}))

	r.Use(middlewares.Func(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, rq *http.Request) {
			val := rq.Context().Value("middleware01").(string)
			g.E(w.Write([]byte(val)))
		})
	}))

	tr := g.Serve()
	tr.Mux.Handle("/", r.ServerHandler())

	g.Eq(g.Req("", tr.URL("/")).String(), "ok")
}

func TestMiddlewareNotFound(t *testing.T) {
	g := got.T(t)

	r := goapi.NewRouter()
	tr := g.Serve()
	tr.Mux.Handle("/", r.ServerHandler())

	r.Use(middlewares.Func(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, rq *http.Request) {
			h.ServeHTTP(w, rq)
		})
	}))

	g.Eq(g.Req("", tr.URL("/x")).StatusCode, http.StatusNotFound)
}

func TestNotFound(t *testing.T) {
	g := got.T(t)

	r := goapi.New()
	tr := g.Serve()
	tr.Mux.Handle("/", r.Server())

	g.Eq(g.Req("", tr.URL("/test")).StatusCode, http.StatusNotFound)
}
