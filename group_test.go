package goapi_test

import (
	"testing"

	"github.com/NaturalSelectionLabs/goapi"
	"github.com/ysmood/got"
)

func TestMultipleGroups(t *testing.T) {
	g := got.T(t)

	r := goapi.New()

	ga := r.Group("/a")
	ga.GET("/users", func() res { return res{Data: "a"} })

	gb := r.Group("/b")
	gb.GET("/users", func() res { return res{Data: "b"} })
	gb.POST("/users", func() res { return res{Data: "b"} })
	gb.PUT("/users", func() res { return res{Data: "b"} })
	gb.PATCH("/users", func() res { return res{Data: "b"} })
	gb.DELETE("/users", func() res { return res{Data: "b"} })
	gb.HEAD("/users", func() res { return res{Data: "b"} })
	gb.OPTIONS("/users", func() res { return res{Data: "b"} })

	tr := g.Serve()
	tr.Mux.Handle("/", r.Server())

	g.Eq(g.Req("", tr.URL("/a/users")).JSON(), map[string]any{
		"data": "a",
	})
	g.Eq(g.Req("", tr.URL("/b/users")).JSON(), map[string]any{
		"data": "b",
	})

	g.Eq(1, 1)
}
