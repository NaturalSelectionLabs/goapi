package goapi

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/iancoleman/strcase"
)

type Router struct {
	middlewares []Middleware
}

func New() *Group {
	r := NewRouter()
	return r.Group("")
}

func NewRouter() *Router {
	return &Router{}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, rq *http.Request) {
	if len(r.middlewares) == 0 {
		writeResponse(w, &ResponseNotFound{
			Error: &Error{
				Message: fmt.Sprintf("path not found: %s %s", rq.Method, rq.URL.Path),
			},
		})

		return
	}

	r.middlewares[0](w, rq, func(w http.ResponseWriter, rq *http.Request) {
		rr := &Router{
			middlewares: r.middlewares[1:],
		}
		rr.ServeHTTP(w, rq)
	})
}

func (r *Router) Add(middleware Middleware) {
	mVal := reflect.ValueOf(middleware)

	if mVal.Kind() != reflect.Func {
		panic("middleware must be a function")
	}

	r.middlewares = append(r.middlewares, middleware)
}

func (r *Router) Group(prefix string) *Group {
	if strcase.ToKebab(prefix) != prefix {
		panic("prefix must be kebab-case: " + prefix)
	}

	g := &Group{
		router:    r,
		prefix:    prefix,
		endpoints: []*Endpoint{},
	}

	return g
}
