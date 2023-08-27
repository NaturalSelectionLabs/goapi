package goapi_test

import (
	"net/http"
	"testing"

	"github.com/NaturalSelectionLabs/goapi"
	"github.com/ysmood/got"
)

type MyResponse goapi.ResponseOK

func (r *MyResponse) StatusCode() int {
	return http.StatusNotModified
}

func TestGroup(t *testing.T) {
	g := got.T(t)

	r := goapi.New()
	tr := g.Serve()
	tr.Mux.Handle("/", r)

	type Params struct {
		ID string
	}

	g.Eq(g.Panic(func() {
		r.GET("/users/{userID}", func() {})
	}), "expect path to be kebab-case, but got: /users/{userID}")

	g.Eq(g.Panic(func() {
		r.GET("/users", "")
	}), "expect handler to be a function, but got: string")

	g.Eq(g.Panic(func() {
		r.GET("/users", func() (int, int, int, int) { return 0, 0, 0, 0 })
	}), "expect handler at most return 3 values, but got: 4")

	r.GET("/users/{id}", func(params *Params) string { return params.ID })

	r.POST("/posts/{id}", func(params *Params) string { return params.ID })

	g.Eq(g.Panic(func() {
		r.GET("/field-type-err", func(params *struct{ P int64 }) {})
	}), "expect struct fields to be string, int, float64, slice of them, or pointer of them, but got: *struct { P int64 }")

	g.Eq(g.Req("", tr.URL("/users/123?user_filter=1")).JSON(), map[string]any{
		"data": "123",
	})

	g.Eq(g.Req(http.MethodPost, tr.URL("/posts/456")).JSON(), map[string]any{
		"data": "456",
	})
}

func TestMultipleGroups(t *testing.T) {
	g := got.T(t)

	r := goapi.NewRouter()
	tr := g.Serve()
	tr.Mux.Handle("/", r)

	ga := r.Group("/a")
	ga.GET("/users", func() string { return "a" })

	gb := r.Group("/b")
	gb.GET("/users", func() string { return "b" })

	g.Eq(g.Req("", tr.URL("/a/users")).JSON(), map[string]any{
		"data": "a",
	})
	g.Eq(g.Req("", tr.URL("/b/users")).JSON(), map[string]any{
		"data": "b",
	})

	g.Eq(1, 1)
}
