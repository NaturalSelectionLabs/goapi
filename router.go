package goapi

import (
	"context"
	"fmt"
	"net/http"

	"github.com/NaturalSelectionLabs/goapi/lib/middlewares"
	"github.com/NaturalSelectionLabs/goapi/lib/openapi"
	"github.com/NaturalSelectionLabs/jschema"
	"github.com/xeipuuv/gojsonschema"
)

// Router itself is a middleware.
type Router struct {
	Schemas jschema.Schemas

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
	s := NewSchemas()

	s.AddTimeHandler()
	s.AddJSONRawMessageHandler()

	return &Router{
		middlewares: []middlewares.Middleware{},
		Schemas:     s,
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
	return middlewares.Chain(r.middlewares...).Handler(next)
}

// Group creates a new group with the given prefix.
func (r *Router) Group(prefix string) *Group {
	g := &Group{router: r}
	return g.Group(prefix)
}

// AddFormatChecker for json schema validation.
// Such as a struct:
//
//	type User struct {
//		ID string `format:"my-id"`
//	}
//
// You can add a format checker for "id" like:
//
//	AddFormatChecker("my-id", checker)
func (r *Router) AddFormatChecker(name string, c gojsonschema.FormatChecker) {
	gojsonschema.FormatCheckers.Add(name, c)
}

func NewSchemas() jschema.Schemas {
	return jschema.NewWithInterfaces("#/components/schemas", interfaces)
}
