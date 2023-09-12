package goapi

import (
	"context"
	"fmt"
	"net/http"

	"github.com/NaturalSelectionLabs/goapi/lib/middlewares"
	"github.com/NaturalSelectionLabs/goapi/lib/openapi"
)

// Router itself is a middleware.
type Router struct {
	// FormatResponse is a function that formats the response.
	// The default format the structs implement [ResponseFormat].
	// You can use this function override the default format.
	FormatResponse FormatResponse
	// Validate the parameters of each request.
	Validate func(v interface{}) *openapi.Error

	middlewares []middlewares.Middleware
	operations  []*Operation
	sever       *http.Server
}

// New is a shortcut for:
//
//	NewRouter().Group("")
func New() *Group {
	return NewRouter().Group("")
}

func NewRouter() *Router {
	return &Router{
		middlewares:    []middlewares.Middleware{},
		FormatResponse: func(format openapi.ResponseFormat) any { return format },
	}
}

// ServerHandler with a 404 middleware at the end.
func (r *Router) ServerHandler() http.Handler {
	return r.Handler(http.HandlerFunc(func(w http.ResponseWriter, rq *http.Request) {
		middlewares.ResponseError(w, http.StatusNotFound, &openapi.Error{
			Code:       openapi.CodeNotFound,
			Message:    fmt.Sprintf("path not found: %s %s", rq.Method, rq.URL.Path),
			Target:     rq.URL.Path,
			InnerError: []any{rq.Method, rq.URL.Path},
		})
	}))
}

// Start listen on addr with the [Router.ServerHandler].
func (r *Router) Start(addr string) error {
	r.sever = &http.Server{
		Addr:    addr,
		Handler: r.ServerHandler(),
	}

	return r.sever.ListenAndServe()
}

// Shutdown the server.
func (r *Router) Shutdown(ctx context.Context) error {
	return r.sever.Shutdown(ctx)
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
