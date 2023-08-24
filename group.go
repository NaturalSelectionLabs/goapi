package goapi

import (
	"net/http"
	"reflect"
	"regexp"
	"strings"

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
	path = g.prefix + path

	if regSpace.MatchString(path) {
		panic("path cannot contain spaces:" + path)
	}

	pathPattern, pathParams, err := openAPIPathToRegexp(path)
	if err != nil {
		panic("path is not a valid OpenAPI path:" + path)
	}

	vHandler := reflect.ValueOf(handler)
	tHandler := vHandler.Type()

	if tHandler.Kind() != reflect.Func {
		panic("handler is not a function")
	}

	if tHandler.NumOut() > 3 {
		panic("handler can at most return 3 values")
	}

	e := &Endpoint{
		method:      method,
		openapiPath: path,
		vHandler:    vHandler,
		tHandler:    tHandler,
		pathPattern: pathPattern,
		pathParams:  pathParams,
	}

	g.endpoints = append([]*Endpoint{e}, g.endpoints...)

	g.router.Add(e.Handle)

	return e
}

// Group is a shortcut for [Router.Group].
func (g *Group) Group(prefix string) *Group {
	return g.router.Group(prefix)
}

// Handler is a shortcut for [Router.Handler].
func (g *Group) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	g.router.ServeHTTP(w, r)
}

var regOpenAPIPath = regexp.MustCompile(`\{([^}]+)\}`)

// Converts OpenAPI style path to Go Regexp and returns path parameters.
func openAPIPathToRegexp(path string) (*regexp.Regexp, []string, error) {
	params := []string{}

	// Replace OpenAPI wildcards with Go RegExp named wildcards
	regexPath := regOpenAPIPath.ReplaceAllStringFunc(path, func(m string) string {
		param := m[1 : len(m)-1] // Strip outer braces from parameter

		if strcase.ToKebab(param) != param { // Make sure parameter is in kebab-case
			panic("path parameter must be in kebab-case: " + param)
		}

		params = append(params, param)    // Add param to list
		return "(?P<" + param + ">[^/]+)" // Replace with Go Regexp named wildcard
	})

	// Make sure the path starts with a "^", ends with a "$", and escape slashes
	regexPath = "^" + strings.ReplaceAll(regexPath, "/", "\\/") + "$"

	// Compile the regular expression
	r, err := regexp.Compile(regexPath)
	if err != nil {
		return nil, nil, err
	}

	return r, params, nil
}
