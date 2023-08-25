package goapi

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/iancoleman/strcase"
)

var (
	tHTTPResponseWriter = reflect.TypeOf((*http.ResponseWriter)(nil)).Elem()
	tHTTPRequest        = reflect.TypeOf((*http.Request)(nil))
	tParamDecoder       = reflect.TypeOf((*ParamDecoder)(nil)).Elem()
)

type Endpoint struct {
	method      string
	openapiPath string
	vHandler    reflect.Value
	tHandler    reflect.Type

	pathPattern    *regexp.Regexp
	pathParams     []string
	overrideWriter bool

	tags []Tag
}

type ParamDecoder interface {
	DecodeParam([]string)
}

func (e *Endpoint) Handle(w http.ResponseWriter, r *http.Request, pathParams []string) {
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

		var field reflect.Value
		if tField.Type.Kind() == reflect.Ptr {
			field = reflect.New(tField.Type.Elem())
		} else {
			field = reflect.New(tField.Type)
		}

		if dec, ok := field.Interface().(ParamDecoder); ok {
			if vs, ok := params[strcase.ToKebab(tField.Name)]; ok {
				dec.DecodeParam(vs)
			} else if vs, ok := qs[strcase.ToSnake(tField.Name)]; ok {
				dec.DecodeParam(vs)
			}

			if tField.Type.Kind() == reflect.Ptr {
				vField.Set(field)
			} else {
				vField.Set(field.Elem())
			}

			continue
		}

		if vs, ok := params[strcase.ToKebab(tField.Name)]; ok {
			assign(tField, vField, vs[0])
		} else if vs, ok := qs[strcase.ToSnake(tField.Name)]; ok {
			if tField.Type.Kind() == reflect.Slice {
				vField.Set(reflect.MakeSlice(tField.Type, len(vs), len(vs)))
				for i, v := range vs {
					assign(tField, vField.Index(i), v)
				}
			} else {
				assign(tField, vField, vs[0])
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

func (e *Endpoint) Tags() []Tag {
	return e.tags
}

type Tag struct {
	name string
}

func (t Tag) String() string {
	return t.name
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

func assign(tField reflect.StructField, vField reflect.Value, val string) {
	kind := vField.Kind()
	ptr := kind == reflect.Ptr

	if ptr {
		kind = tField.Type.Elem().Kind()
	}

	switch kind { //nolint: exhaustive
	case reflect.Int:
		n, _ := strconv.ParseInt(val, 10, 64)
		setVal(vField, ptr, reflect.ValueOf(int(n)))

	case reflect.Float64:
		n, _ := strconv.ParseFloat(val, 64)
		setVal(vField, ptr, reflect.ValueOf(n))

	default:
		setVal(vField, ptr, reflect.ValueOf(val))
	}
}

func setVal(vField reflect.Value, ptr bool, val reflect.Value) {
	if ptr {
		p := reflect.New(val.Type())
		p.Elem().Set(val)
		val = p
	}

	vField.Set(val)
}
