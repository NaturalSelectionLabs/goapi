package goapi

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/iancoleman/strcase"
)

type Middleware interface {
	// A Middleware https://cs.opensource.google/go/x/pkgsite/+/68be0dd1:internal/middleware/middleware.go
	Handler(http.Handler) http.Handler
}

type MiddlewareFunc func(http.Handler) http.Handler

func (fn MiddlewareFunc) Handler(h http.Handler) http.Handler {
	return fn(h)
}

// Router itself is a middleware.
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

func NewRouter() *Router {
	return &Router{}
}

func (r *Router) Server() http.Handler {
	return r.Handler(http.HandlerFunc(func(w http.ResponseWriter, rq *http.Request) {
		writeResErr(w, http.StatusNotFound, fmt.Sprintf("path not found: %s %s", rq.Method, rq.URL.Path))
	}))
}

// Add a middleware to the router.
func (r *Router) Add(middlewares ...Middleware) {
	r.middlewares = append(r.middlewares, middlewares...)
}

func (r *Router) Handler(h http.Handler) http.Handler {
	for i := len(r.middlewares) - 1; i >= 0; i-- {
		h = r.middlewares[i].Handler(h)
	}

	return h
}

// Group creates a new group with the given prefix.
func (r *Router) Group(prefix string) *Group {
	if len(prefix) > 0 && prefix[len(prefix)-1] == '/' {
		panic("expect prefix to not end with '/', but got: " + prefix)
	}

	if strcase.ToKebab(prefix) != prefix {
		panic("expect prefix be kebab-cased, but got: " + prefix)
	}

	if regexp.MustCompile(`[{}]`).MatchString(prefix) {
		panic("expect prefix not contains braces, but got: " + prefix)
	}

	g := &Group{
		router: r,
		prefix: prefix,
	}

	return g
}
