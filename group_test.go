package goapi_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/NaturalSelectionLabs/goapi"
	"github.com/ysmood/got"
)

func TestMultipleGroups(t *testing.T) {
	g := got.T(t)

	r := goapi.New()

	{
		ga := r.Group("/a")
		ga.GET("/users", func() res { return resOK{Data: "a"} })
		gb := r.Group("/b")

		gb.GET("/users", func() res { return resOK{Data: "b"} })
		gb.POST("/users", func() res { return resOK{Data: "post"} })
		gb.PUT("/users", func() res { return resOK{Data: "b"} })
		gb.PATCH("/users", func() res { return resOK{Data: "b"} })
		gb.DELETE("/users", func() res { return resOK{Data: "b"} })
		gb.HEAD("/users", func() res { return resOK{Data: "b"} })
		gb.OPTIONS("/users", func() res { return resOK{Data: "b"} })

		g.Eq(g.Panic(func() {
			gb.Group("user")
		}), "expect prefix to start with '/', but got: user")

		g.Eq(g.Panic(func() {
			gb.Group("/user/")
		}), "expect prefix to not end with '/', but got: /user/")

		g.Eq(g.Panic(func() {
			gb.Group("/us_er")
		}), "expect prefix be kebab-cased, but got: /us_er")

		g.Eq(g.Panic(func() {
			gb.Group("/{user}")
		}), "expect prefix not contains braces, but got: /{user}")
	}

	tr := g.Serve()
	tr.Mux.Handle("/", r.Server())

	g.Eq(g.Req("", tr.URL("/a/users")).JSON(), map[string]any{
		"data": "a",
	})
	g.Eq(g.Req("", tr.URL("/b/users")).JSON(), map[string]any{
		"data": "b",
	})
	g.Eq(g.Req(http.MethodPost, tr.URL("/b/users")).JSON(), map[string]any{
		"data": "post",
	})

	g.Eq(1, 1)
}

func TestStart(t *testing.T) {
	g := got.T(t)
	r := goapi.New()
	r.GET("/", func() resOK { return resOK{Data: "ok"} })

	go func() { _ = r.Start(":41375") }()

	time.Sleep(300 * time.Millisecond)

	g.Eq(g.Req("", "http://localhost:41375").JSON(), map[string]interface{}{
		"data": "ok",
	})

	g.E(r.Shutdown(g.Context()))
}
