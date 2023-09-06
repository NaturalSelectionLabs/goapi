package goapi

import (
	"fmt"
	"net/http"
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
	return &Router{
		middlewares: []Middleware{
			MiddlewareFunc(func(h http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, rq *http.Request) {
					defer func() {
						if err := recover(); err != nil {
							writeResErr(w, http.StatusInternalServerError, fmt.Sprint(err))
						}
					}()

					h.ServeHTTP(w, rq)
				})
			}),
		},
	}
}

func (r *Router) Server() http.Handler {
	return r.Handler(http.HandlerFunc(func(w http.ResponseWriter, rq *http.Request) {
		writeResErr(w, http.StatusNotFound, fmt.Sprintf("path not found: %s %s", rq.Method, rq.URL.Path))
	}))
}

// Use a middleware to the router.
func (r *Router) Use(middlewares ...Middleware) {
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
	g := &Group{
		router: r,
		prefix: prefix,
	}

	return g
}
