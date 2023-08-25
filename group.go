package goapi

import (
	"fmt"
	"net/http"
	"reflect"
	"regexp"

	"github.com/iancoleman/strcase"
)

type Group struct {
	router    *Router
	prefix    string
	endpoints []*Endpoint
}

var regSpace = regexp.MustCompile(`\s+`)

func (g *Group) GET(path string, handler any) *Endpoint {
	return g.Add(http.MethodGet, path, handler)
}

func (g *Group) POST(path string, handler any) *Endpoint {
	return g.Add(http.MethodPost, path, handler)
}

func (g *Group) PUT(path string, handler any) *Endpoint {
	return g.Add(http.MethodPut, path, handler)
}

func (g *Group) PATCH(path string, handler any) *Endpoint {
	return g.Add(http.MethodPatch, path, handler)
}

func (g *Group) DELETE(path string, handler any) *Endpoint {
	return g.Add(http.MethodDelete, path, handler)
}

func (g *Group) OPTIONS(path string, handler any) *Endpoint {
	return g.Add(http.MethodOptions, path, handler)
}

func (g *Group) HEAD(path string, handler any) *Endpoint {
	return g.Add(http.MethodHead, path, handler)
}

func (g *Group) Add(method, path string, handler any) *Endpoint {
	if strcase.ToKebab(path) != path {
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

	e := &Endpoint{
		method:      method,
		openapiPath: path,
		vHandler:    vHandler,
		tHandler:    tHandler,
		pathPattern: pathPattern,
		pathParams:  pathParams,
	}

	for i := 0; i < tHandler.NumIn(); i++ {
		tArg := tHandler.In(i)

		switch tArg {
		case tHTTPResponseWriter:
			e.overrideWriter = true

		case tHTTPRequest:

		default:
			if tArg.Kind() != reflect.Ptr || tArg.Elem().Kind() != reflect.Struct {
				panic("expect handler arguments must be http.ResponseWriter, *http.Request, or pointer to a struct, " +
					"but got: " + tArg.String())
			}

			for j := 0; j < tArg.Elem().NumField(); j++ {
				t := tArg.Elem().Field(j).Type
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

	g.endpoints = append([]*Endpoint{e}, g.endpoints...)

	return e
}

// Group creates a sub group of current group.
func (g *Group) Group(prefix string) *Group {
	return g.router.Group(g.prefix + prefix)
}

// Handler is a shortcut for [Router.Handler].
func (g *Group) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	g.router.ServeHTTP(w, r)
}
