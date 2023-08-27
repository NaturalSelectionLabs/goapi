package goapi

import (
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"strings"

	"github.com/NaturalSelectionLabs/goapi/lib/openapi"
)

type Group struct {
	router     *Router
	prefix     string
	operations []*Operation
}

var regSpace = regexp.MustCompile(`\s+`)

func (g *Group) GET(path string, handler any, opts ...GroupAddOption) *Operation {
	return g.Add(openapi.GET, path, handler, opts...)
}

func (g *Group) POST(path string, handler any, opts ...GroupAddOption) *Operation {
	return g.Add(openapi.POST, path, handler, opts...)
}

func (g *Group) PUT(path string, handler any, opts ...GroupAddOption) *Operation {
	return g.Add(openapi.PUT, path, handler, opts...)
}

func (g *Group) PATCH(path string, handler any, opts ...GroupAddOption) *Operation {
	return g.Add(openapi.PATCH, path, handler, opts...)
}

func (g *Group) DELETE(path string, handler any, opts ...GroupAddOption) *Operation {
	return g.Add(openapi.DELETE, path, handler, opts...)
}

func (g *Group) OPTIONS(path string, handler any, opts ...GroupAddOption) *Operation {
	return g.Add(openapi.OPTIONS, path, handler, opts...)
}

func (g *Group) HEAD(path string, handler any, opts ...GroupAddOption) *Operation {
	return g.Add(openapi.HEAD, path, handler, opts...)
}

type GroupAddOption func(op *Operation)

// Meta is a type of option for [Group.Add] to set the meta info of an operation.
func (g *Group) Meta(meta OperationMeta) GroupAddOption {
	return func(op *Operation) { op.meta = &meta }
}

func (g *Group) Add( //nolint: gocognit
	method openapi.Method, path string, handler any, opts ...GroupAddOption,
) *Operation {
	if toPathName(path) != path {
		panic("expect path to be kebab-case, but got: " + path)
	}

	if regSpace.MatchString(path) {
		panic("expect path contain no spaces, but got: " + path)
	}

	pathPattern, pathParams, err := openAPIPathToRegexp(path)
	if err != nil {
		panic("expect path matches the openapi path format, but got: " + path)
	}

	vHandler := reflect.ValueOf(handler)
	tHandler := vHandler.Type()

	if tHandler.Kind() != reflect.Func {
		panic("expect handler to be a function, but got: " + tHandler.String())
	}

	if tHandler.NumOut() > 3 {
		panic(fmt.Sprintf("expect handler at most return 3 values, but got: %d", tHandler.NumOut()))
	}

	op := &Operation{
		method:      method,
		openapiPath: path,
		vHandler:    vHandler,
		tHandler:    tHandler,
		pathPattern: pathPattern,
		pathParams:  pathParams,
	}

	for _, opt := range opts {
		opt(op)
	}

	for i := 0; i < tHandler.NumIn(); i++ {
		tArg := tHandler.In(i)

		switch tArg {
		case tHTTPResponseWriter:
			op.overrideWriter = true

		case tHTTPRequest:

		default:
			if tArg.Kind() != reflect.Ptr || tArg.Elem().Kind() != reflect.Struct {
				panic("expect handler arguments must be http.ResponseWriter, *http.Request, or pointer to a struct, " +
					"but got: " + tArg.String())
			}

			for j := 0; j < tArg.Elem().NumField(); j++ {
				t := tArg.Elem().Field(j).Type

				if t.Implements(tParamDecoder) || reflect.New(t).Type().Implements(tParamDecoder) {
					continue
				}

				if t.Kind() == reflect.Ptr {
					t = t.Elem()
				}

				if t.Kind() == reflect.Slice {
					t = t.Elem()
				}

				switch t.Kind() { //nolint: exhaustive
				case reflect.String, reflect.Int, reflect.Float64:
				default:
					panic("expect struct fields to be string, int, " +
						"float64, slice of them, or pointer of them, but got: " + tArg.String())
				}
			}
		}
	}

	g.operations = append([]*Operation{op}, g.operations...)

	return op
}

func (g *Group) Handle(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	path := r.URL.Path

	if !strings.HasPrefix(path, g.prefix) {
		next(w, r)
		return
	}

	for _, op := range g.operations {
		pathParams := op.Match(r.Method, path[len(g.prefix):])
		if pathParams != nil {
			op.Handle(w, r, pathParams)
			return
		}
	}

	next(w, r)
}

// Group creates a sub group of current group.
func (g *Group) Group(prefix string) *Group {
	return g.router.Group(g.prefix + prefix)
}

// Handler is a shortcut for [Router.Handler].
func (g *Group) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	g.router.ServeHTTP(w, r)
}
