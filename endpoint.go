package goapi

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"strings"

	"github.com/iancoleman/strcase"
)

var (
	tHTTPResponseWriter = reflect.TypeOf((*http.ResponseWriter)(nil)).Elem()
	tHTTPRequest        = reflect.TypeOf((*http.Request)(nil))
)

type Endpoint struct {
	method      string
	openapiPath string
	vHandler    reflect.Value
	tHandler    reflect.Type
	tags        []Tag

	pathPattern    *regexp.Regexp
	pathParams     []string
	overrideWriter bool
}

func (e *Endpoint) Handle(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	pathParams := e.Match(r.Method, r.URL.Path)

	if pathParams == nil {
		next(w, r)
		return
	}

	args := make([]reflect.Value, e.tHandler.NumIn())

	for i := range args {
		tArg := e.tHandler.In(i)

		switch tArg {
		case tHTTPResponseWriter:
			args[i] = reflect.ValueOf(w)

		case tHTTPRequest:
			args[i] = reflect.ValueOf(r)

		default:
			qs := r.URL.Query()

			if err := e.guardQuery(qs); err != nil {
				writeResponse(w, &ResponseBadRequest{Error: toError(err)})
				return
			}

			args[i] = e.Params(tArg.Elem(), pathParams, qs)
		}
	}

	ret := e.vHandler.Call(args)

	if e.overrideWriter {
		return
	}

	if len(ret) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if len(ret) > 0 {
		last := ret[len(ret)-1].Interface()

		if err, ok := last.(error); ok {
			if _, ok := last.(ErrorCode); ok {
				writeResponse(w, &ResponseBadRequest{Error: toError(err)})
			} else {
				writeResponse(w, &ResponseInternalServerError{Error: toError(err)})
			}

			return
		}
	}

	first := ret[0].Interface()

	if res, ok := first.(Response); ok {
		writeResponse(w, res)
		return
	}

	switch len(ret) {
	case 1:
		writeResponse(w, &ResponseOK{first, nil})

	default:
		writeResponse(w, &ResponseOK{first, ret[1].Interface()})
	}
}

func (e *Endpoint) Match(method, path string) []string {
	if method != e.method {
		return nil
	}

	matches := e.pathPattern.FindStringSubmatch(path)

	if matches == nil {
		return nil
	}

	return matches[1:]
}

func (e *Endpoint) Params(tArg reflect.Type, pathParams []string, qs url.Values) reflect.Value {
	arg := reflect.New(tArg)

	params := make(map[string][]string, len(pathParams))

	for i, name := range e.pathParams {
		params[name] = []string{pathParams[i]}
	}

	for i := 0; i < tArg.NumField(); i++ {
		tField := tArg.Field(i)
		vField := arg.Elem().Field(i)

		if v, ok := params[strcase.ToKebab(tField.Name)]; ok {
			vField.SetString(v[0])
		} else if v, ok := qs[strcase.ToSnake(tField.Name)]; ok {
			if tField.Type.Kind() == reflect.Slice {
				vField.Set(reflect.ValueOf(v))
			} else {
				vField.SetString(v[0])
			}
		}
	}

	return arg
}

// Ensure query key are camelCased.
func (e *Endpoint) guardQuery(qs url.Values) error {
	for k := range qs {
		if k != strcase.ToSnake(k) {
			return fmt.Errorf("query key is not snake styled: %s", k)
		}
	}

	return nil
}

func (e *Endpoint) SetTags(tags ...Tag) *Endpoint {
	e.tags = tags
	return e
}

type Tag struct {
	name string
}

var regOpenAPIPath = regexp.MustCompile(`\{([^}]+)\}`)

// Converts OpenAPI style path to Go Regexp and returns path parameters.
func openAPIPathToRegexp(path string) (*regexp.Regexp, []string, error) {
	params := []string{}

	// Replace OpenAPI wildcards with Go RegExp named wildcards
	regexPath := regOpenAPIPath.ReplaceAllStringFunc(path, func(m string) string {
		param := m[1 : len(m)-1]          // Strip outer braces from parameter
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
