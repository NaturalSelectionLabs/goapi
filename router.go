package goapi

import (
	"context"
	"fmt"
	"net/http"
)

type Middleware interface {
	// A Middleware https://cs.opensource.google/go/x/pkgsite/+/68be0dd1:internal/middleware/middleware.go
	Handler(next http.Handler) http.Handler
}

type MiddlewareFunc func(next http.Handler) http.Handler

func (fn MiddlewareFunc) Handler(next http.Handler) http.Handler {
	return fn(next)
}

// Router itself is a middleware.
type Router struct {
	middlewares []Middleware
	operations  []*Operation
	Server      *http.Server
}

// New is a shortcut for:
//
//	NewRouter().Group("")
func New() *Group {
	return NewRouter().Group("")
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

// ServerHandler with a 404 middleware at the end.
func (r *Router) ServerHandler() http.Handler {
	return r.Handler(http.HandlerFunc(func(w http.ResponseWriter, rq *http.Request) {
		writeResErr(w, http.StatusNotFound, fmt.Sprintf("path not found: %s %s", rq.Method, rq.URL.Path))
	}))
}

// Start the server.
func (r *Router) Start(addr string) error {
	r.Server = &http.Server{
		Addr:    addr,
		Handler: r.ServerHandler(),
	}

	return r.Server.ListenAndServe()
}

// Shutdown the server.
func (r *Router) Shutdown(ctx context.Context) error {
	return r.Server.Shutdown(ctx)
}

// Use a middleware to the router.
func (r *Router) Use(middlewares ...Middleware) {
	r.middlewares = append(r.middlewares, middlewares...)
}

func (r *Router) Handler(next http.Handler) http.Handler {
	for i := len(r.middlewares) - 1; i >= 0; i-- {
		next = r.middlewares[i].Handler(next)
	}

	return next
}

// Group creates a new group with the given prefix.
func (r *Router) Group(prefix string) *Group {
	g := &Group{router: r}
	return g.Group(prefix)
}
