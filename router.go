package goapi

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/iancoleman/strcase"
)

type Router struct {
	middlewares []Middleware
}

// New is a shortcut for:
//
//	NewRouter().Group("")
func New() *Group {
	r := NewRouter()
	return r.Group("")
}

// NewRouter creates a new router, it implements [http.Handler].
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

	r.middlewares[0].Handle(w, rq, func(w http.ResponseWriter, rq *http.Request) {
		rr := &Router{
			middlewares: r.middlewares[1:],
		}
		rr.ServeHTTP(w, rq)
	})
}

// Add a middleware to the router.
func (r *Router) Add(middleware Middleware) {
	r.middlewares = append(r.middlewares, middleware)
}

var regBraces = regexp.MustCompile(`[{}]`)

// Group creates a new group with the given prefix.
func (r *Router) Group(prefix string) *Group {
	if strcase.ToKebab(prefix) != prefix {
		panic("expect prefix be kebab-cased, but got: " + prefix)
	}

	if regBraces.MatchString(prefix) {
		panic("expect prefix not contains braces, but got: " + prefix)
	}

	g := &Group{
		router:     r,
		prefix:     prefix,
		operations: []*Operation{},
	}

	r.Add(g)

	return g
}
