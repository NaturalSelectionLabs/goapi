package goapi

import (
	"context"
	"fmt"
	"net/http"

	"github.com/NaturalSelectionLabs/goapi/lib/middlewares"
)

// Router itself is a middleware.
type Router struct {
	middlewares []middlewares.Middleware
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
		middlewares: []middlewares.Middleware{},
	}
}

// ServerHandler with a 404 middleware at the end.
func (r *Router) ServerHandler() http.Handler {
	return r.Handler(http.HandlerFunc(func(w http.ResponseWriter, rq *http.Request) {
		middlewares.ResponseError(w, http.StatusNotFound, fmt.Sprintf("path not found: %s %s", rq.Method, rq.URL.Path))
	}))
}

// Start listen on addr with the [Router.ServerHandler].
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
func (r *Router) Use(middlewares ...middlewares.Middleware) {
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
