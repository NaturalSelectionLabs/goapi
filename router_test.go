package goapi_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/NaturalSelectionLabs/goapi"
	"github.com/NaturalSelectionLabs/jschema"
	"github.com/ysmood/got"
)

var s = jschema.New("")

type Activity interface {
	Summary() string
}

var IActivity = jschema.DefineI(s, new(Activity))

type Post struct {
	Title string
	Body  string
}

var _ = IActivity.Define(Post{})

func (p Post) Summary() string {
	return p.Title
}

type Transaction struct {
	From   string
	To     string
	Amount float64
}

var _ = IActivity.Define(Transaction{})

func (t Transaction) Summary() string {
	return fmt.Sprintf("%s -> %s: %f", t.From, t.To, t.Amount)
}

type Filter struct {
	Keyword string
	Limit   int
	Offset  int
}

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
