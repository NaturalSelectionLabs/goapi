package goapi

import (
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"strings"

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

var regBraces = regexp.MustCompile(`[{}]`)

func (r *Router) Group(prefix string) *Group {
	if strcase.ToKebab(prefix) != prefix {
		panic("expect prefix be kebab-cased, but got: " + prefix)
	}

	if regBraces.MatchString(prefix) {
		panic("expect prefix not contains braces, but got: " + prefix)
	}

	g := &Group{
		router:    r,
		prefix:    prefix,
		endpoints: []*Endpoint{},
	}

	r.Add(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		path := r.URL.Path

		if !strings.HasPrefix(path, prefix) {
			next(w, r)
			return
		}

		for _, e := range g.endpoints {
			pathParams := e.Match(r.Method, path[len(prefix):])
			if pathParams != nil {
				e.Handle(w, r, pathParams)
				return
			}
		}
	})

	return g
}
