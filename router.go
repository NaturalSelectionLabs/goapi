package goapi

import (
	"net/http"
	"reflect"
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
		w.WriteHeader(http.StatusNotFound)
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
	g := &Group{
		router:    r,
		prefix:    prefix,
		endpoints: []*Endpoint{},
	}

	return g
}
