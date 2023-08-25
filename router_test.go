package goapi_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/NaturalSelectionLabs/goapi"
	"github.com/ysmood/got"
)

func TestMiddleware(t *testing.T) {
	g := got.T(t)

	r := goapi.NewRouter()

	r.Add(func(w http.ResponseWriter, rq *http.Request, next http.HandlerFunc) {
		rq = rq.WithContext(context.WithValue(rq.Context(), "middleware01", "ok")) //nolint: staticcheck
		next(w, rq)
	})
	r.Add(func(w http.ResponseWriter, rq *http.Request, _ http.HandlerFunc) {
		val := rq.Context().Value("middleware01").(string)
		g.E(w.Write([]byte(val)))
	})

	tr := g.Serve()
	tr.Mux.Handle("/", r)

	g.Eq(g.Req("", tr.URL("/")).String(), "ok")
}

func TestGroupErr(t *testing.T) {
	g := got.T(t)

	r := goapi.NewRouter()

	g.Eq(g.Panic(func() {
		r.Group("/users/{id}")
	}), "expect prefix not contains braces, but got: /users/{id}")

	g.Eq(1, 1)
}
