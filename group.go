package goapi

import (
	"context"
	"net/http"
	"reflect"
	"regexp"
	"strings"

	ff "github.com/NaturalSelectionLabs/goapi/lib/flat-fields"
	"github.com/NaturalSelectionLabs/goapi/lib/middlewares"
	"github.com/NaturalSelectionLabs/goapi/lib/openapi"
	"github.com/xeipuuv/gojsonschema"
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

// GET is a shortcut for [Group.Add].
func (g *Group) GET(path string, handler OperationHandler) {
	g.Add(openapi.GET, path, handler)
}

// POST is a shortcut for [Group.Add].
func (g *Group) POST(path string, handler OperationHandler) {
	g.Add(openapi.POST, path, handler)
}

// PUT is a shortcut for [Group.Add].
func (g *Group) PUT(path string, handler OperationHandler) {
	g.Add(openapi.PUT, path, handler)
}

// PATCH is a shortcut for [Group.Add].
func (g *Group) PATCH(path string, handler OperationHandler) {
	g.Add(openapi.PATCH, path, handler)
}

// DELETE is a shortcut for [Group.Add].
func (g *Group) DELETE(path string, handler OperationHandler) {
	g.Add(openapi.DELETE, path, handler)
}

// OPTIONS is a shortcut for [Group.Add].
func (g *Group) OPTIONS(path string, handler OperationHandler) {
	g.Add(openapi.OPTIONS, path, handler)
}

// HEAD is a shortcut for [Group.Add].
func (g *Group) HEAD(path string, handler OperationHandler) {
	g.Add(openapi.HEAD, path, handler)
}

// Add adds a new http handler to the group.
// If a request matches the path and method, the handler will be called.
// The router will ignore the trailing slash of the path if a path without trailing slash
// has not been defined.
func (g *Group) Add(method openapi.Method, path string, handler OperationHandler) {
	op := g.newOperation(method, g.prefix+path, handler)
	g.router.operations = append(g.router.operations, op)
	g.Use(op)
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

func (g *Group) parseParam(path *Path, p reflect.Type) *parsedParam {
	if p == tContext {
		return &parsedParam{isContext: true}
	}

	if p == tRequest {
		return &parsedParam{isRequest: true}
	}

	if p.Kind() != reflect.Struct || !p.Implements(tParams) {
		panic("expect parameter to be a struct and embedded with goapi.InHeader, goapi.InURL, or goapi.InBody," +
			" but got: " + p.String())
	}

	parsed := &parsedParam{param: p}
	fields := []*parsedField{}
	flat := ff.Parse(p)

	switch reflect.New(p).Interface().(Params).paramsIn() {
	case inHeader:
		parsed.in = inHeader

		for _, f := range flat.Fields {
			fields = append(fields, g.parseHeaderField(f))
		}

	case inURL:
		parsed.in = inURL

		for _, f := range flat.Fields {
			fields = append(fields, g.parseURLField(path, f))
		}

		for _, n := range path.names {
			has := false

			for _, f := range fields {
				if f.name == n {
					has = true
				}
			}

			if !has {
				panic("expect to have path parameter for {" + n + "} in " + p.String())
			}
		}

	case inBody:
		parsed.in = inBody

		scm := g.router.Schemas.ToStandAlone(g.router.Schemas.DefineT(p))

		validator, _ := gojsonschema.NewSchema(gojsonschema.NewGoLoader(scm))
		parsed.bodyValidator = validator
	}

	parsed.fields = fields

	return parsed
}

func (g *Group) parseHeaderField(flatField *ff.FlattenedField) *parsedField {
	f := flatField.Field
	parsed := g.parseField(flatField)
	parsed.name = toHeaderName(f.Name)
	parsed.name = tagName(f.Tag, parsed.name)

	return parsed
}

func (g *Group) parseURLField(path *Path, flatField *ff.FlattenedField) *parsedField {
	f := g.parseField(flatField)

	t := flatField.Field

	f.name = toPathName(t.Name)
	if path.contains(f.name) {
		if f.hasDefault {
			panic("path parameter cannot have tag `default`, param: " + t.Name)
		}

		if f.slice {
			panic("path parameter cannot be an slice, param: " + t.Name)
		}

		if !f.required {
			panic("path parameter cannot be optional, param: " + t.Name)
		}

		f.InPath = true
	} else {
		f.name = toQueryName(t.Name)
	}

	f.name = tagName(t.Tag, f.name)

	return f
}

func (g *Group) parseField(flatField *ff.FlattenedField) *parsedField {
	t := flatField.Field
	f := &parsedField{flatField: flatField, required: true}
	tf := t.Type

	switch t.Type.Kind() { //nolint: exhaustive
	case reflect.Ptr, reflect.Slice:
		f.ptr = true
	default:
		f.ptr = false
	}

	if tf.Kind() == reflect.Ptr {
		tf = tf.Elem()
		f.required = false
	}

	if tf.Kind() == reflect.Slice {
		f.slice = true
		f.sliceType = tf
		f.item = tf.Elem()
		f.required = false
	} else {
		f.item = tf
	}

	f.schema = g.router.Schemas.ToStandAlone(firstProp(g.router.Schemas.DefineFieldT(t)))

	if _, ok := t.Tag.Lookup("default"); ok {
		f.required = false
		f.hasDefault = true
		f.defaultVal = reflect.ValueOf(f.schema.Default)
	}

	if _, ok := t.Tag.Lookup("example"); ok {
		f.example = reflect.ValueOf(f.schema.Example)
	}

	validator, _ := gojsonschema.NewSchema(gojsonschema.NewGoLoader(f.schema))
	f.validator = validator

	return f
}
