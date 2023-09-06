package goapi

import (
	"net/http"
	"regexp"

	"github.com/NaturalSelectionLabs/goapi/lib/openapi"
	"github.com/iancoleman/strcase"
)

type Group struct {
	router *Router
	prefix string
}

func (g *Group) GET(path string, handler any, opts ...ConfigOperation) {
	g.Add(openapi.GET, path, handler, opts...)
}

func (g *Group) POST(path string, handler any, opts ...ConfigOperation) {
	g.Add(openapi.POST, path, handler, opts...)
}

func (g *Group) PUT(path string, handler any, opts ...ConfigOperation) {
	g.Add(openapi.PUT, path, handler, opts...)
}

func (g *Group) PATCH(path string, handler any, opts ...ConfigOperation) {
	g.Add(openapi.PATCH, path, handler, opts...)
}

func (g *Group) DELETE(path string, handler any, opts ...ConfigOperation) {
	g.Add(openapi.DELETE, path, handler, opts...)
}

func (g *Group) OPTIONS(path string, handler any, opts ...ConfigOperation) {
	g.Add(openapi.OPTIONS, path, handler, opts...)
}

func (g *Group) HEAD(path string, handler any, opts ...ConfigOperation) {
	g.Add(openapi.HEAD, path, handler, opts...)
}

type ConfigOperation func(op *Operation)

// Meta is a type of option for [Group.Add] to set the meta info of an operation.
func (g *Group) Meta(meta OperationMeta) ConfigOperation {
	return func(op *Operation) { op.meta = &meta }
}

func (g *Group) Add(
	method openapi.Method, path string, handler any, opts ...ConfigOperation,
) {
	op := newOperation(method, g.prefix+path, handler)

	for _, opt := range opts {
		opt(op)
	}

	g.router.Use(op)
}

// Group creates a sub group of current group.
func (g *Group) Group(prefix string) *Group {
	if len(prefix) > 0 && prefix[0] != '/' {
		panic("expect prefix to start with '/', but got: " + prefix)
	}

	if len(prefix) > 0 && prefix[len(prefix)-1] == '/' {
		panic("expect prefix to not end with '/', but got: " + prefix)
	}

	if strcase.ToKebab(prefix) != prefix {
		panic("expect prefix be kebab-cased, but got: " + prefix)
	}

	if regexp.MustCompile(`[{}]`).MatchString(prefix) {
		panic("expect prefix not contains braces, but got: " + prefix)
	}

	return g.router.Group(g.prefix + prefix)
}

// Handler is a shortcut for [Router.Handler].
func (g *Group) Server() http.Handler {
	return g.router.Server()
}

// Use is a shortcut for [Router.Use].
func (g *Group) Use(m Middleware) {
	g.router.Use(m)
}
