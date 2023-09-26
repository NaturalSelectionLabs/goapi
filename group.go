package goapi

import (
	"context"
	"net/http"
	"regexp"
	"strings"

	"github.com/NaturalSelectionLabs/goapi/lib/middlewares"
	"github.com/NaturalSelectionLabs/goapi/lib/openapi"
)

// Group of handlers with the same url path prefix.
type Group struct {
	router *Router
	prefix string
}

// Router returns the router of the group.
func (g *Group) Router() *Router {
	return g.router
}

// Prefix returns the prefix of the group.
func (g *Group) Prefix() string {
	return g.prefix
}

// Add is a shortcut for [Group.Add] for a random golang function.
// It only accepts POST request method.
// It will use the function name as the url path name.
// It will treat the input and output of the function as the request and response body.
func Add[P, S any](g *Group, fn func(P) S) *Operation {
	name := toPathName(fnName(fn))

	type res struct {
		StatusOK
		Data S `response:"direct"`
	}

	return g.POST("/"+name, func(p P) res {
		return res{Data: fn(p)}
	})
}

// GET is a shortcut for [Group.Add].
func (g *Group) GET(path string, handler OperationHandler) *Operation {
	return g.Add(openapi.GET, path, handler)
}

// POST is a shortcut for [Group.Add].
func (g *Group) POST(path string, handler OperationHandler) *Operation {
	return g.Add(openapi.POST, path, handler)
}

// PUT is a shortcut for [Group.Add].
func (g *Group) PUT(path string, handler OperationHandler) *Operation {
	return g.Add(openapi.PUT, path, handler)
}

// PATCH is a shortcut for [Group.Add].
func (g *Group) PATCH(path string, handler OperationHandler) *Operation {
	return g.Add(openapi.PATCH, path, handler)
}

// DELETE is a shortcut for [Group.Add].
func (g *Group) DELETE(path string, handler OperationHandler) *Operation {
	return g.Add(openapi.DELETE, path, handler)
}

// OPTIONS is a shortcut for [Group.Add].
func (g *Group) OPTIONS(path string, handler OperationHandler) *Operation {
	return g.Add(openapi.OPTIONS, path, handler)
}

// HEAD is a shortcut for [Group.Add].
func (g *Group) HEAD(path string, handler OperationHandler) *Operation {
	return g.Add(openapi.HEAD, path, handler)
}

// Add adds a new http handler to the group.
// If a request matches the path and method, the handler will be called.
// The router will ignore the trailing slash of the path if a path without trailing slash
// has not been defined.
func (g *Group) Add(method openapi.Method, path string, handler OperationHandler) *Operation {
	op := g.newOperation(method, g.prefix+path, handler)
	g.router.operations = append(g.router.operations, op)
	g.Use(op)

	return op
}

// Group creates a sub group of current group.
func (g *Group) Group(prefix string) *Group {
	if len(prefix) > 0 && prefix[0] != '/' {
		panic("expect prefix to start with '/', but got: " + prefix)
	}

	if len(prefix) > 0 && prefix[len(prefix)-1] == '/' {
		panic("expect prefix to not end with '/', but got: " + prefix)
	}

	if regexp.MustCompile(`[{}]`).MatchString(prefix) {
		panic("expect prefix not contains braces, but got: " + prefix)
	}

	return &Group{
		router: g.router,
		prefix: g.prefix + prefix,
	}
}

// Server is a shortcut for [Router.Handler].
func (g *Group) Server() http.Handler {
	return g.router.ServerHandler()
}

// Start is a shortcut for [Router.Start].
func (g *Group) Start(addr string) error {
	return g.router.Start(addr)
}

// Shutdown is a shortcut for [Router.Shutdown].
func (g *Group) Shutdown(ctx context.Context) error {
	return g.router.Shutdown(ctx)
}

// Use is similar to [Router.Use] but with he group prefix.
func (g *Group) Use(m middlewares.Middleware) {
	g.router.Use(middlewares.Func(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, g.prefix) {
				m.Handler(h).ServeHTTP(w, r)
			} else {
				h.ServeHTTP(w, r)
			}
		})
	}))
}

// Handler is a shortcut for [Router.Handler].
func (g *Group) Handler(h http.Handler) http.Handler {
	return g.router.Handler(h)
}
